package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	feederLock sync.Mutex

	latestUpdatedTime int64
)

func (app *App) StartDataAutomaticFeeder() error {
	config := getTimeConfig()
	app.logger.Infof("Data Automatic Feeder is starting with interval %v and will immediatly feed when longer than %v", config.interval, config.maximumDelay)

	tickers := []*time.Ticker{
		schedule(app.startFeeding, config.interval),
		schedule(app.forceFeeding, config.maximumDelay),
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

func (app *App) startFeeding() {
	app.ctx = context.Background()
	feederLock.Lock()
	defer feederLock.Unlock()
	app.logger.Infof("feeding process is starting")

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
	priceMapping, err := app.getRequestedPricingFromSource(reqId, config)
	if err != nil {
		app.logger.Errorf("could not get requested pricing from source because: %v", err)
		return
	}

	// prepare data before request update latest price information to destination
	updateParams := priceMapping.toDstParams()
	app.logger.Debugf("%v", updateParams)
	// TODO
}

func (app *App) forceFeeding() {
	app.ctx = context.Background()

	feederLock.Lock()
	defer feederLock.Unlock()
	app.logger.Debugf("checking if needed to force update pricing to destination")

}
