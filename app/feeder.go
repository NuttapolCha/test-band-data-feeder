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
	logger := app.logger

	config := getTimeConfig()
	logger.Infof("Data Automatic Feeder is starting")

	tickers := []*time.Ticker{
		schedule(app.getDataFromSource, config.dataSourceInterval),
		schedule(app.updateDataToDestination, config.destinationInterval),
	}

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	logger.Infof("Data Automatic Feeder has stopped")

	for _, ticker := range tickers {
		ticker.Stop()
	}
	return nil
}

func (app *App) getDataFromSource() {
	logger := app.logger

	dataSourceLock.Lock()
	defer dataSourceLock.Unlock()
	logger.Infof("getting data from source..")

	defer func() {
		if r := recover(); r != nil {
			logger.Debugf("panic and recover because: %v", r)
		}
	}()

	config := getFeederConfig()
	logger.Debugf("symbols: %v", config.symbols)

	// request pricing information from data source
	reqId, err := app.requestPricingFromSource(config)
	if err != nil {
		logger.Errorf("could not request pricing from source because: %v", err)
		return
	}

	// some delay before getting the requested priceing
	time.Sleep(config.waitTime * time.Second)

	// get pricing data from the requested
	pricingResults, err := app.getRequestedPricingFromSource(reqId, config)
	if err != nil {
		logger.Errorf("could not get requested pricing from source because: %v", err)
		return
	}

	// caching accquired pricing results to memory
	err = app.cachePricingResults(pricingResults)
	if err != nil {
		logger.Errorf("could not cache pricing results to memory because: %v", err)
		return
	}
}

func (app *App) updateDataToDestination() {
	logger := app.logger

	destinationLock.Lock()
	defer destinationLock.Unlock()
	logger.Debugf("checking if we needed to update pricing to destination..")

	defer func() {
		if r := recover(); r != nil {
			logger.Debugf("panic and recover because: %v", r)
		}
	}()

	config := getFeederConfig()

	// mustUpdatePricingList describes pricing that needed to send update at destination
	// because it meet 1 of 2 conditions.
	// 1. no new update than 1 hour (configurable)
	// 2. pricing difference is more than threshold 0.1 (configurable)
	mustUpdatePricingList := make([]*UpdatePricingParams, 0)

	type pair struct {
		symbol string
		price  float64
	}
	timestampMapSymbolPricing := make(map[int64][]pair)

	for _, symbol := range config.symbols {
		// get latest price from cache
		latestPriceFromCache, err := cache.GetLatestPrice(symbol)
		if err != nil {
			logger.Errorf("could not get latest price of %s from cache because: %v", symbol, err)
			return
		}
		logger.Debugf("%s last update at %d price = %.4f", symbol, latestPriceFromCache.GetTimestamp(), latestPriceFromCache.GetPrice())

		// get latest price from destination
		latestPriceFromDst, err := app.getLatestPriceFromDestination(symbol, config)
		if err != nil {
			logger.Warnf("could not get latest price of %s from destinationn because: %v", symbol, err)
			logger.Warnf("destination maybe never received %s information before and need to send update", symbol)
			latestPriceFromDst = new(DestinationPricingResp)
		}

		// prepare data for update classified by timestamp
		if app.isNeedUpdatePricingToDestination(symbol, latestPriceFromCache, latestPriceFromDst, config) {
			latestTimestamp := latestPriceFromCache.GetTimestamp()
			timestampMapSymbolPricing[latestTimestamp] = append(
				timestampMapSymbolPricing[latestTimestamp],
				pair{
					symbol: symbol,
					price:  latestPriceFromCache.GetPrice(),
				},
			)
		}
	}

	// append to mustUpdatePricingList
	for latestTimestamp, symbolPricePairs := range timestampMapSymbolPricing {
		// separate each fields to 2 slices
		prices := make([]float64, 0, len(symbolPricePairs))
		symbols := make([]string, 0, len(symbolPricePairs))
		for _, symbolPricePair := range symbolPricePairs {
			prices = append(prices, symbolPricePair.price)
			symbols = append(symbols, symbolPricePair.symbol)
		}

		mustUpdatePricingList = append(mustUpdatePricingList, &UpdatePricingParams{
			Symbols:   symbols,
			Prices:    prices,
			Timestamp: latestTimestamp,
		})
	}

	// update pricing to destination
	updatedSymbols, err := app.updatePricingToDestination(mustUpdatePricingList, config)
	if err != nil {
		logger.Errorf("unsuccessful update pricing to destination because: %v", err)
	}
	logger.Infof("updated symbols are %+v", updatedSymbols)
}
