package cmd

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"ftp-server/internal/config"
)

var (
	version string
)

var rootCmd = &cobra.Command{
	Use:   "ftp-server",
	Short: "Ftp Server",
	RunE:  run,
}

func init() {
	cobra.OnInitialize(initConfig)
	viper.SetDefault("ftp.user", "ds")
	viper.SetDefault("ftp.password", "testwalk.123")
	viper.SetDefault("ftp.port", 2222)
	viper.SetDefault("ftp.host", "localhost")
	viper.SetDefault("ftp.passive-port", "2223-2224")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)
}

// Execute executes the root command.
func Execute(v string) {
	version = v

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func initConfig() {
	config.Version = version
	viper.SetConfigName("ftp-server")
	viper.AddConfigPath("$GOPATH/src/ftp-server/config")
	viper.AddConfigPath("./")
	viper.AddConfigPath("/etc/ftp-server")
	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Warning("No configuration file found, using defaults.")
		default:
			log.WithError(err).Fatal("read configuration file error")
		}
	}
	for _, pair := range os.Environ() {
		d := strings.SplitN(pair, "=", 2)
		if strings.Contains(d[0], ".") {
			log.Warning("Using dots in env variable is illegal and deprecated. Please use double underscore `__` for: ", d[0])
			underscoreName := strings.ReplaceAll(d[0], ".", "__")
			// Set only when the underscore version doesn't already exist.
			if _, exists := os.LookupEnv(underscoreName); !exists {
				os.Setenv(underscoreName, d[1])
			}
		}
	}
	viperBindEnvs(config.C)
	viperHooks := mapstructure.ComposeDecodeHookFunc(
		viperDecodeJSONSlice,
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	)
	if err := viper.Unmarshal(&config.C, viper.DecodeHook(viperHooks)); err != nil {
		log.WithError(err).Fatal("unmarshal config error")
	}

}

func viperDecodeJSONSlice(rf reflect.Kind, rt reflect.Kind, data interface{}) (interface{}, error) {
	// input must be a string and destination must be a slice
	if rf != reflect.String || rt != reflect.Slice {
		return data, nil
	}

	raw := data.(string)

	// this decoder expects a JSON list
	if !strings.HasPrefix(raw, "[") || !strings.HasSuffix(raw, "]") {
		return data, nil
	}

	var out []map[string]interface{}
	err := json.Unmarshal([]byte(raw), &out)

	return out, err
}

func viperBindEnvs(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			tv = strings.ToLower(t.Name)
		}
		if tv == "-" {
			continue
		}

		switch v.Kind() {
		case reflect.Struct:
			viperBindEnvs(v.Interface(), append(parts, tv)...)
		default:
			// Bash doesn't allow env variable names with a dot so
			// bind the double underscore version.
			keyDot := strings.Join(append(parts, tv), ".")
			keyUnderscore := strings.Join(append(parts, tv), "__")
			viper.BindEnv(keyDot, strings.ToUpper(keyUnderscore))
		}
	}
}
