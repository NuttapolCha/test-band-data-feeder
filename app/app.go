package app

import (
	"context"

	"github.com/NuttapolCha/test-band-data-feeder/log"
)

type App struct {
	logger log.Logger
	ctx    context.Context
}

func New(logger log.Logger, ctx context.Context) App {
	return App{
		logger: logger,
		ctx:    ctx,
	}
}

func (app *App) DataAutomaticFeeder() error {
	return nil
}
