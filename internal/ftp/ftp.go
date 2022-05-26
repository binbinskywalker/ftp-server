package ftp

import (
	"ftp-server/internal/config"
	"log"
	"os"

	filedriver "github.com/goftp/file-driver"
	"github.com/goftp/server"
)

var ftpServer *server.Server

var (
	user     string
	password string
	port     int
	host     string
	group    string
	owner    string
	dataport string
)

const (
	rootPath = "/usr/local/ftp-server"
)

// Setup ftp server
func Setup(c config.Config) error {
	user = c.Ftp.User
	password = c.Ftp.PassWord
	port = c.Ftp.Port
	host = c.Ftp.Host
	group = c.Ftp.Group
	owner = c.Ftp.Owner
	dataport = c.Ftp.DataPort

	if _, err := PathExists(rootPath); err != nil {
		return err
	}

	factory := &filedriver.FileDriverFactory{
		RootPath: rootPath,
		Perm:     server.NewSimplePerm(owner, group),
	}

	opts := &server.ServerOpts{
		Factory:      factory,
		Port:         port,
		PassivePorts: dataport,
		Hostname:     host,
		Name:         "ds-ftp-server",
		Auth:         &server.SimpleAuth{Name: user, Password: password},
	}

	go func(opts *server.ServerOpts) {
		ftpServer = server.NewServer(opts)
		err := ftpServer.ListenAndServe()
		if err != nil {
			log.Fatal("Error starting server:", err)
		}
	}(opts)

	return nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	return false, err
}

// Stop stops ftp handler.
func Stop() error {
	return ftpServer.Shutdown()
}
