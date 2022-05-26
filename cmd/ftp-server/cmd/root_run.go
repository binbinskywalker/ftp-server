package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"ftp-server/internal/config"
	"ftp-server/internal/ftp"
)

func run(cmd *cobra.Command, args []string) error {
	if err := ftp.Setup(config.C); err != nil {
		return errors.Wrap(err, "setup ftp server error")
	}

	sigChan := make(chan os.Signal)
	exitChan := make(chan struct{})
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	log.WithField("signal", <-sigChan).Info("signal received")

	go func() {
		if err := ftp.Stop(); err != nil {
			fmt.Println(err)
		}
		exitChan <- struct{}{}
	}()
	select {
	case <-exitChan:
	case s := <-sigChan:
		log.WithField("signal", s).Info("signal received, stopping immediately")
	}
	return nil
}
