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

// StartDataAutomaticFeeder called by cmd after initialized application.
// It does get data from source, caching in memory and update to destination if neccessary
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

// Feed called by cmd and run only once.
func (app *App) Feed() {
	logger := app.logger
	logger.Infof("Feed is starting")
	app.getDataAndFeed()
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
	symbolMapPricing := make(map[string]pricing.Information)

	// check for each symbol need to update to destination or not
	for _, currPricing := range pricingResults {
		var is, immediatly bool
		var prevUpdateDstTime int64
		var err error

		symbol := currPricing.Symbol
		prevPricing, err := cache.GetPricing(symbol)
		if err != nil {
			logger.Infof("no previous pricing information found in cache, need update to destination")
			is = true
			goto sendToDestination
		}
		prevUpdateDstTime, err = cache.GetPrevUpdatedDstTime(symbol)
		if err != nil {
			logger.Errorf("could not get previous updated destination time of %s because: %v", symbol, err)
			return
		}

		is, immediatly = app.isNeedUpdatePricingToDestination(prevUpdateDstTime, prevPricing, currPricing, config)
		if immediatly {
			// force update this symbol now (we cannot wait)
			done := make(chan struct{})
			go func(symbol string, price float64) {
				urgentMap := map[string]pricing.Information{
					symbol: currPricing,
				}
				updatedSymbol, err := app.updatePricingToDestination(urgentMap, config)
				if err != nil {
					logger.Errorf("could not update %s pricing to destination immediatly because: %v", symbol, err)
					return
				}
				logger.Infof("successfully updated %+v pricing to destination immediatly", updatedSymbol)

				done <- struct{}{}
			}(currPricing.GetSymbol(), currPricing.GetPrice())

			<-done
			continue
		}

	sendToDestination:
		if is {
			symbolMapPricing[symbol] = currPricing
		}
	}

	// update pricing to destination
	updatedSymbols, err := app.updatePricingToDestination(symbolMapPricing, config)
	if err != nil {
		logger.Errorf("update pricing to destination not completed because: %v", err)
		return
	}
	logger.Infof("updated symbols for this interval (exclude immediatly sent) are %+v", updatedSymbols)
}
