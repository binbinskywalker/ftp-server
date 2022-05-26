package cmd

import (
	"os"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"ftp-server/internal/config"
)

// when updating this template, don't forget to update config.md!
const configTemplate = `
  # Ftp settings.
  [ftp]
  # string
  user="{{ .Ftp.User }}"
  # string
  password="{{ .Ftp.PassWord }}"
  # int
  port={{ .Ftp.Port }}
  # string
  host="{{ .Ftp.Host }}"
  group="{{ .Ftp.Group }}"
  owner="{{ .Ftp.Owner }}"
  passive-port="{{ .Ftp.DataPort }}"
`

var configCmd = &cobra.Command{
	Use:   "configfile",
	Short: "Print the ChirpStack Network Server configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		t := template.Must(template.New("config").Parse(configTemplate))
		err := t.Execute(os.Stdout, &config.C)
		if err != nil {
			return errors.Wrap(err, "execute config template error")
		}
		return nil
	},
}
