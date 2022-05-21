package app

import (
	"time"

	"github.com/spf13/viper"
)

type TimeConfig struct {
	dataSourceInterval  time.Duration
	destinationInterval time.Duration
}

var timeConfig *TimeConfig

func getTimeConfig() *TimeConfig {
	if timeConfig == nil {
		timeConfig = &TimeConfig{
			dataSourceInterval:  viper.GetDuration("DataFeeder.DataSourceInterval") * time.Second,
			destinationInterval: viper.GetDuration("DataFeeder.DestinationInterval") * time.Second,
		}
	}
	return timeConfig
}
