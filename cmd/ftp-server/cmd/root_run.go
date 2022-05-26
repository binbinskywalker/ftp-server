package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"ftp-server/internal/config"
	"ftp-server/internal/ftp"
)

func run(cmd *cobra.Command, args []string) error {
	if err := ftp.Setup(config.C); err != nil {
		return errors.Wrap(err, "setup ftp server error")
	}
	return nil
}
