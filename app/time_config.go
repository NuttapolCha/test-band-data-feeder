package app

import (
	"time"

	"github.com/spf13/viper"
)

type TimeConfig struct {
	interval time.Duration
}

var timeConfig *TimeConfig

func getTimeConfig() *TimeConfig {
	if timeConfig == nil {
		timeConfig = &TimeConfig{
			interval: viper.GetDuration("DataFeeder.Interval") * time.Second,
		}
		if timeConfig.interval == 0 {
			timeConfig.interval = 10 * time.Second
		}
	}
	return timeConfig
}
