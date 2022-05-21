package app

import (
	"time"

	"github.com/spf13/viper"
)

type FeederConfig struct {
	diffThreshold float64

	// Symbols
	symbols []string

	// Data Source
	dataSourceRetryCount       int
	requestPricingDataEndpoint string
	getPricingDataEndpoint     string

	// wait time before get the requested pricing
	waitTime time.Duration

	// Update Destination
	destinationRetryCount     int
	updatePricingDataEndpoint string
	getUpdatedPricingData     string
}

var feederConfig *FeederConfig

func getFeederConfig() *FeederConfig {
	if feederConfig == nil {
		feederConfig = &FeederConfig{
			symbols:                    viper.GetStringSlice("DataFeeder.Symbols"),
			waitTime:                   viper.GetDuration("DataFeeder.WaitTime"),
			diffThreshold:              viper.GetFloat64("DataFeeder.DiffThreshold"),
			dataSourceRetryCount:       viper.GetInt("ExternalAPIs.DataSource.RetryCount"),
			requestPricingDataEndpoint: viper.GetString("ExternalAPIs.DataSource.RequestPricingData"),
			getPricingDataEndpoint:     viper.GetString("ExternalAPIs.DataSource.GetPricingData"),
			destinationRetryCount:      viper.GetInt("ExternalAPIs.Destination.RetryCount"),
			updatePricingDataEndpoint:  viper.GetString("ExternalAPIs.Destination.UpdatePricingData"),
			getUpdatedPricingData:      viper.GetString("ExternalAPIs.Destination.GetUpdatedPricingData"),
		}
	}
	return feederConfig
}