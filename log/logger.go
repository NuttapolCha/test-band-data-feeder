package log

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type logLevel int

const (
	info logLevel = iota
	debug
	verbose
)

type Logger struct {
	level logLevel
}

// NewLogger initializes a simple logger
func NewLogger() (Logger, error) {
	var lvl logLevel

	switch strings.ToLower(viper.GetString("Log.Level")) {
	case "debug":
		lvl = debug
	case "verbose":
		lvl = verbose
	default:
		lvl = info
	}

	logger := Logger{
		level: lvl,
	}

	return logger, nil
}

func (logger *Logger) Infof(template string, args ...interface{}) {
	log.Printf("\033[0;34m[INFO]\033[0;37m "+template+"\n", args...)
}

func (logger *Logger) Errorf(template string, args ...interface{}) {
	log.Printf("\033[0;31m[ERROR]\033[0;37m "+template+"\n", args...)
}

func (logger *Logger) Warnf(template string, args ...interface{}) {
	log.Printf("\033[0;33m[WARN]\033[0;37m "+template+"\n", args...)
}

func (logger *Logger) Debugf(template string, args ...interface{}) {
	if logger.level < debug {
		return
	}
	log.Printf("\033[0;35m[DEBUG]\033[0;37m "+template+"\n", args...)
}

func (logger *Logger) BeautyJSON(bs []byte) {
	if logger.level < verbose {
		return
	}

	var i interface{}
	json.Unmarshal(bs, &i)

	res, _ := json.MarshalIndent(&i, "", "\t")
	log.Printf("\033[0;35m[DEBUG]\033[0;37m " + string(res) + "\n")
}
