package app

import (
	"time"

	"github.com/spf13/viper"
)

type FeederConfig struct {
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
	maximumDelay              int64
	diffThreshold             float64
	enableRecheck             bool
}

var feederConfig *FeederConfig

func getFeederConfig() *FeederConfig {
	if feederConfig == nil {
		feederConfig = &FeederConfig{
			symbols:                    viper.GetStringSlice("DataFeeder.Symbols"),
			waitTime:                   viper.GetDuration("DataFeeder.WaitTime") * time.Second,
			dataSourceRetryCount:       viper.GetInt("ExternalAPIs.DataSource.RetryCount"),
			requestPricingDataEndpoint: viper.GetString("ExternalAPIs.DataSource.RequestPricingData"),
			getPricingDataEndpoint:     viper.GetString("ExternalAPIs.DataSource.GetPricingData"),
			destinationRetryCount:      viper.GetInt("ExternalAPIs.Destination.RetryCount"),
			updatePricingDataEndpoint:  viper.GetString("ExternalAPIs.Destination.UpdatePricingData"),
			getUpdatedPricingData:      viper.GetString("ExternalAPIs.Destination.GetUpdatedPricingData"),
			maximumDelay:               viper.GetInt64("DataFeeder.MaximumDelay"),
			diffThreshold:              viper.GetFloat64("DataFeeder.DiffThreshold"),
			enableRecheck:              viper.GetBool("DataFeeder.EnableRecheck"),
		}
		if feederConfig.symbols == nil {
			feederConfig.symbols = []string{"BTC", "ETH"}
		}
		if feederConfig.waitTime == 0 {
			feederConfig.waitTime = 5 * time.Second
		}
		if feederConfig.dataSourceRetryCount == 0 {
			feederConfig.dataSourceRetryCount = 1
		}
		if feederConfig.requestPricingDataEndpoint == "" {
			feederConfig.requestPricingDataEndpoint = "https://interview-requester-source.herokuapp.com/request"
		}
		if feederConfig.getPricingDataEndpoint == "" {
			feederConfig.getPricingDataEndpoint = "https://interview-requester-source.herokuapp.com/request"
		}
		if feederConfig.destinationRetryCount == 0 {
			feederConfig.destinationRetryCount = 1
		}
		if feederConfig.updatePricingDataEndpoint == "" {
			feederConfig.updatePricingDataEndpoint = "https://band-interview-destination.herokuapp.com/update"
		}
		if feederConfig.getUpdatedPricingData == "" {
			feederConfig.getUpdatedPricingData = "https://band-interview-destination.herokuapp.com/get_price"
		}
		if feederConfig.maximumDelay == 0 {
			feederConfig.maximumDelay = 3600
		}
		if feederConfig.diffThreshold == 0 {
			feederConfig.diffThreshold = 0.1
		}

	}
	return feederConfig
}
