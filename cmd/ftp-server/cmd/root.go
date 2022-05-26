package cmd

import (
	"reflect"
	"strings"

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
	viperBindEnvs(config.C)
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
