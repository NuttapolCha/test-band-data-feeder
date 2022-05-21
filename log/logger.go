package log

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type logLevel int

const (
	info logLevel = iota
	debug
)

type Logger struct{ level logLevel }

// NewLogger initializes a simple logger
func NewLogger() Logger {
	return Logger{
		level: func() logLevel {
			if strings.ToLower(viper.GetString("Log.Level")) == "debug" {
				return debug
			}
			return info
		}(),
	}
}

func (logger *Logger) Infof(template string, args ...interface{}) {
	log.Printf("\033[0;34m[INFO]\033[0;37m "+template+"\n", args...)
}

func (logger *Logger) Debugf(template string, args ...interface{}) {
	if logger.level != debug {
		return
	}
	log.Printf("\033[0;35m[DEBUG]\033[0;37m "+template+"\n", args...)
}

func (logger *Logger) Errorf(template string, args ...interface{}) {
	log.Printf("\033[0;31m[ERROR]\033[0;37m "+template+"\n", args...)
}

func (logger *Logger) Warnf(template string, args ...interface{}) {
	log.Printf("\033[0;33m[WARN]\033[0;37m "+template+"\n", args...)
}
