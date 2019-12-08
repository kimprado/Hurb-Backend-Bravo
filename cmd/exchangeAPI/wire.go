// +build wireinject

package main

import (
	"github.com/google/wire"
	app "github.com/rep/exchange/internal/app"
	"github.com/rep/exchange/internal/pkg/commom/config"
	"github.com/rep/exchange/internal/pkg/commom/logging"
)

// initializeConfig inicializa Configuration
func initializeConfig(path string) (config config.Configuration, err error) {
	panic(wire.Build(app.AppSet))
}

// initializeAppender inicializa FileAppender
func initializeAppender(path string) (appender logging.FileAppender, err error) {
	panic(wire.Build(app.AppSet))
}

// initializeApp inicializa ExchangeApp
func initializeApp(path string) (a *app.ExchangeApp, err error) {
	panic(wire.Build(app.AppSet))
}
