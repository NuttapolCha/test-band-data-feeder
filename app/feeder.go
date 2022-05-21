package app

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/NuttapolCha/test-band-data-feeder/cache"
)

var (
	dataSourceLock  sync.Mutex
	destinationLock sync.Mutex
)

func (app *App) StartDataAutomaticFeeder() error {
	config := getTimeConfig()
	app.logger.Infof("Data Automatic Feeder is starting")

	tickers := []*time.Ticker{
		schedule(app.getDataFromSource, config.dataSourceInterval),
		schedule(app.updateDataToDestination, config.destinationInterval),
	}

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	app.logger.Infof("Data Automatic Feeder has stopped")

	for _, ticker := range tickers {
		ticker.Stop()
	}
	return nil
}

func (app *App) getDataFromSource() {
	dataSourceLock.Lock()
	defer dataSourceLock.Unlock()
	app.logger.Infof("getting data from source..")

	config := getFeederConfig()
	app.logger.Debugf("symbols: %v", config.symbols)

	// request pricing information from data source
	reqId, err := app.requestPricingFromSource(config)
	if err != nil {
		app.logger.Errorf("could not request pricing from source because: %v", err)
		return
	}

	// some delay before getting the requested priceing
	time.Sleep(config.waitTime * time.Second)

	// get pricing data from the requested
	pricingResults, err := app.getRequestedPricingFromSource(reqId, config)
	if err != nil {
		app.logger.Errorf("could not get requested pricing from source because: %v", err)
		return
	}

	// caching accquired pricing results to memory
	err = app.cachePricingResults(pricingResults)
	if err != nil {
		app.logger.Errorf("could not cache pricing results to memory because: %v", err)
		return
	}
}

func (app *App) updateDataToDestination() {
	destinationLock.Lock()
	defer destinationLock.Unlock()
	app.logger.Debugf("checking if we needed to update pricing to destination..")

	config := getFeederConfig()

	for _, sym := range config.symbols {
		latestPrice, err := cache.GetLatestPrice(sym)
		if err != nil {
			app.logger.Errorf("could not get latest price of %s from cache because: %v", sym, err)
			return
		}
		app.logger.Infof("%s last update at %d price = %.4f", sym, latestPrice.Timestamp, latestPrice.Price)
	}

	// TODO
}
