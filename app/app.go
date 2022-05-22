package app

import (
	"context"
	"time"

	"github.com/NuttapolCha/test-band-data-feeder/connector"
	"github.com/NuttapolCha/test-band-data-feeder/log"
)

type App struct {
	logger     log.Logger
	ctx        context.Context
	httpClient *connector.CustomHttpClient
}

func New(logger log.Logger, httpClient *connector.CustomHttpClient) App {
	return App{
		logger:     logger,
		ctx:        context.TODO(),
		httpClient: httpClient,
	}
}

func schedule(f func(), d time.Duration) *time.Ticker {
	ticker := time.NewTicker(d)
	go func() {
		for range ticker.C {
			f()
		}
	}()
	return ticker
}
