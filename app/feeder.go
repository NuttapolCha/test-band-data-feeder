package app

import (
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	"github.com/NuttapolCha/test-band-data-feeder/app/pricing"
	"github.com/NuttapolCha/test-band-data-feeder/cache"
)

var (
	dataSourceLock sync.Mutex
)

func (app *App) StartDataAutomaticFeeder() error {
	logger := app.logger

	config := getTimeConfig()
	logger.Infof("Data Automatic Feeder is starting")

	tickers := []*time.Ticker{
		schedule(app.getDataAndFeed, config.interval),
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

func (app *App) getDataAndFeed() {
	logger := app.logger

	dataSourceLock.Lock()
	defer dataSourceLock.Unlock()
	logger.Infof("getting data from source..")

	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("panic and recover because: %v", r)
			logger.Debugf("debug stack = %s", debug.Stack())
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

	// declare variables related to updating pricing to destination
	timestampMapPricing := make(map[int64][]pricing.Information)
	mustUpdatePricingList := make([]*UpdatePricingParams, 0)

	// check for each symbol need to update to destination or not
	for _, currPricing := range pricingResults {
		var is, immediatly bool

		symbol := currPricing.Symbol
		prevPricing, err := cache.GetPricing(symbol)
		if err != nil {
			logger.Infof("no previous pricing information found in cache, results in need update to destination")
			is = true
			goto sendToDestination
		}

		is, immediatly = app.isNeedUpdatePricingToDestination(prevPricing, currPricing, config)
		if immediatly {
			// force update this symbol now (we cannot wait)
			go func(symbol string, price float64) {
				updatedSymbol, err := app.updatePricingToDestination([]*UpdatePricingParams{
					{
						Symbols:   []string{symbol},
						Prices:    []float64{price},
						Timestamp: currPricing.GetTimestamp(),
					},
				}, config)
				if err != nil {
					logger.Errorf("could not update %s pricing to destination immediatly because: %v", symbol, err)
					return
				}
				logger.Infof("successfully updated %+v pricing to destination immediatly", updatedSymbol)

				// TODO: might refactor this section because function called is duplicated below
				// cache new current pricing after retreived previous pricing
				if err := cache.UpdatePricing(symbol, currPricing.GetPrice()); err != nil {
					logger.Warnf("could not cache current pricing of %s because: %v", symbol, err)
				}
			}(currPricing.GetSymbol(), currPricing.GetPrice())

			continue
		}

	sendToDestination:
		if is {
			timestampMapPricing[currPricing.GetTimestamp()] = append(
				timestampMapPricing[currPricing.GetTimestamp()],
				currPricing,
			)
		}

		// cache new current pricing after retreived previous pricing
		if err := cache.UpdatePricing(symbol, currPricing.GetPrice()); err != nil {
			logger.Warnf("could not cache current pricing of %s because: %v", symbol, err)
		}
	}

	// append to mustUpdatePricingList
	for t, pricingList := range timestampMapPricing {
		// separate each fields to 2 slices
		prices := make([]float64, 0, len(pricingList))
		symbols := make([]string, 0, len(pricingList))
		for _, pricingInfo := range pricingList {
			prices = append(prices, pricingInfo.GetPrice())
			symbols = append(symbols, pricingInfo.GetSymbol())
		}

		mustUpdatePricingList = append(mustUpdatePricingList, &UpdatePricingParams{
			Symbols:   symbols,
			Prices:    prices,
			Timestamp: t,
		})
	}

	// update pricing to destination
	updatedSymbols, err := app.updatePricingToDestination(mustUpdatePricingList, config)
	if err != nil {
		logger.Errorf("update pricing to destination not completed because: %v", err)
		return
	}
	logger.Infof("updated symbols for this interval are %+v", updatedSymbols)

	// TODO: recheck destination by query its latest pricing
}
