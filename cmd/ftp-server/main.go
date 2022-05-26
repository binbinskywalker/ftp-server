package main

import (
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/grpclog"

	"ftp-server/cmd/ftp-server/cmd"
)

type grpcLogger struct {
	*log.Logger
}

func (gl *grpcLogger) V(l int) bool {
	level, ok := map[log.Level]int{
		log.DebugLevel: 0,
		log.InfoLevel:  1,
		log.WarnLevel:  2,
		log.ErrorLevel: 3,
		log.FatalLevel: 4,
	}[log.GetLevel()]
	if !ok {
		return false
	}

	return l >= level
}

func (gl *grpcLogger) Info(args ...interface{}) {
	if log.GetLevel() == log.DebugLevel {
		log.Debug(args...)
	}
}

func (gl *grpcLogger) Infoln(args ...interface{}) {
	if log.GetLevel() == log.DebugLevel {
		log.Debug(args...)
	}
}

func (gl *grpcLogger) Infof(format string, args ...interface{}) {
	if log.GetLevel() == log.DebugLevel {
		log.Debugf(format, args...)
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: time.RFC3339Nano,
	})

	grpclog.SetLoggerV2(&grpcLogger{log.StandardLogger()})
}

var version string // set by the compiler

func main() {
	cmd.Execute(version)
}
