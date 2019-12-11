// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/rep/exchange/internal/app"
	"github.com/rep/exchange/internal/pkg/commom/config"
	"github.com/rep/exchange/internal/pkg/commom/logging"
	"github.com/rep/exchange/internal/pkg/currencyexchange"
	"github.com/rep/exchange/internal/pkg/currencyexchange/api"
	"github.com/rep/exchange/internal/pkg/infra/redis"
	"github.com/rep/exchange/internal/pkg/webserver"
)

// Injectors from wire.go:

func initializeConfig(path string) (config.Configuration, error) {
	configuration, err := config.NewConfig(path)
	if err != nil {
		return config.Configuration{}, err
	}
	return configuration, nil
}

func initializeAppender(path string) (logging.FileAppender, error) {
	configuration, err := config.NewConfig(path)
	if err != nil {
		return logging.FileAppender{}, err
	}
	fileAppender := logging.NewFileAppender(configuration)
	return fileAppender, nil
}

func initializeApp(path string) (*app.ExchangeApp, error) {
	configuration, err := config.NewConfig(path)
	if err != nil {
		return nil, err
	}
	redisDB := config.NewRedisDB(configuration)
	loggingLevels := config.NewLoggingLevels(configuration)
	loggerRedisDB := logging.NewRedisDB(loggingLevels)
	dbConnection, err := redis.NewDBConnection(redisDB, loggerRedisDB)
	if err != nil {
		return nil, err
	}
	loggerCurrency := logging.NewCurrency(loggingLevels)
	currencyManagerDB := currencyexchange.NewCurrencyManagerDB(dbConnection, redisDB, loggerCurrency)
	currencyManagerProxy := currencyexchange.NewCurrencyManagerProxy(currencyManagerDB, loggerCurrency)
	loggerRates := logging.NewRates(loggingLevels)
	ratesFinderService := currencyexchange.NewRatesFinderService(configuration, loggerRates)
	loggerRatesCache := logging.NewRatesCache(loggingLevels)
	ratesFinderCache := currencyexchange.NewRatesFinderCache(ratesFinderService, configuration, loggerRatesCache)
	ratesFinderProxy := currencyexchange.NewRatesFinderProxy(ratesFinderCache, loggerRates)
	currencyExchanger := currencyexchange.NewCurrencyExchanger()
	loggerCalculator := logging.NewCalculator(loggingLevels)
	calculatorController := currencyexchange.NewCalculatorController(currencyManagerProxy, ratesFinderProxy, currencyExchanger, loggerCalculator)
	loggerAPIExchange := logging.NewLoggerAPIExchange(loggingLevels)
	controller := api.NewController(calculatorController, loggerAPIExchange)
	loggerWebServer := logging.NewWebServer(loggingLevels)
	paramWebServer := webserver.NewParamWebServer(controller, configuration, loggerWebServer)
	webServer := webserver.NewWebServer(paramWebServer)
	logger := logging.NewLogger(loggingLevels)
	exchangeApp := app.NewExchangeApp(webServer, logger)
	return exchangeApp, nil
}
