package cache

import (
	"sync"
	"time"

	"github.com/NuttapolCha/test-band-data-feeder/log"
	"github.com/NuttapolCha/test-band-data-feeder/utils"
	"github.com/spf13/viper"
)

type Pricing struct {
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

func (p *Pricing) GetPrice() float64 {
	return p.Price
}

func (p *Pricing) GetTimestamp() int64 {
	return p.Timestamp
}

type SymbolMapPricing map[string]*Pricing

type pricingSyncMap struct {
	mu sync.Mutex
	m  SymbolMapPricing
}

var (
	latestDataSource  pricingSyncMap
	latestDestination pricingSyncMap

	liveTime                  time.Duration
	destinationCacheUpdatedAt time.Time
	destinationCachePath      string
	dataSourceCachePath       string
)

// Init initialize data source and destination cache
func Init(logger log.Logger) {
	liveTime = viper.GetDuration("DataFeeder.Cache.LiveTime") * time.Second
	destinationCachePath = viper.GetString("DataFeeder.Cache.Destination")
	dataSourceCachePath = viper.GetString("DataFeeder.Cache.DataSource")

	err := initDataSource(logger)
	if err != nil {
		logger.Errorf("intit data source cache got unexpected result because: %v", err)
	}
	err = initDestination(logger)
	if err != nil {
		logger.Errorf("intit destination cache got unexpected result because: %v", err)
	}
}

// Done update cache JSON file using latest application memory
func Done(logger log.Logger) {
	if err := utils.CreateFile(dataSourceCachePath, latestDataSource.m); err != nil {
		logger.Errorf("could not create data source cache file because: %v", err)
	} else {
		logger.Infof("data source cache has been inited")
	}
	if err := utils.CreateFile(destinationCachePath, latestDestination.m); err != nil {
		logger.Errorf("could not create destination cache file because: %v", err)
	} else {
		logger.Infof("destination cache has been inited")
	}
}

func initDataSource(logger log.Logger) error {
	latestDataSource = pricingSyncMap{
		mu: sync.Mutex{},
		m:  make(SymbolMapPricing),
	}

	if !utils.IsFileOld(dataSourceCachePath, liveTime) {
		logger.Infof("data source cache is not old, continue using it")
		err := utils.UnmarshalFromFile(dataSourceCachePath, &latestDataSource.m)
		if err != nil {
			return err
		}
	} else {
		logger.Infof("data source cache is old, create a new one")
	}
	return utils.CreateFile(dataSourceCachePath, &latestDataSource.m)
}

func initDestination(logger log.Logger) error {
	latestDestination = pricingSyncMap{
		mu: sync.Mutex{},
		m:  make(SymbolMapPricing),
	}

	if !utils.IsFileOld(destinationCachePath, liveTime) {
		logger.Infof("data source cache is not old, continue using it")
		err := utils.UnmarshalFromFile(destinationCachePath, &latestDestination.m)
		if err != nil {
			return err
		}
	} else {
		logger.Infof("destination cache is old, create a new one")
	}
	return utils.CreateFile(destinationCachePath, latestDestination.m)
}
