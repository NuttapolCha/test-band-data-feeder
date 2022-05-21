package app

import (
	"time"

	"github.com/spf13/viper"
)

type TimeConfig struct {
	interval     time.Duration
	maximumDelay time.Duration
}

var timeConfig *TimeConfig

func getTimeConfig() *TimeConfig {
	if timeConfig == nil {
		timeConfig = &TimeConfig{
			interval:     viper.GetDuration("DataFeeder.TriggeredInterval") * time.Second,
			maximumDelay: viper.GetDuration("DataFeeder.MaximumDelay") * time.Second,
		}
	}
	return timeConfig
}
