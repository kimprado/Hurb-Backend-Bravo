// +build wireinject

package currencyexchange

import (
	"github.com/google/wire"
	"github.com/rep/exchange/internal/pkg/commom/config"
	"github.com/rep/exchange/internal/pkg/infra/redis"
)

// initializeConfigTest inicializa Configuration para testes
func initializeConfigTest() (config config.Configuration, err error) {
	panic(wire.Build(pkgSetConfigTest))
}

// initializeRedisTest inicializa DBConnection para testes
func initializeRedisTest(config config.Configuration) (c redis.DBConnection, err error) {
	panic(wire.Build(pkgSetTest))
}

// initializeRedisTest inicializa CurrencyManagerDB para testes
func initializeCurrencyManagerDBTest(config config.Configuration) (c *CurrencyManagerDB, err error) {
	panic(wire.Build(pkgSetTest))
}

// initializeRatesFinderService inicializa RatesFinderService para testes
func initializeRatesFinderServiceTest(config config.Configuration) (c *RatesFinderService, err error) {
	panic(wire.Build(pkgSetTest))
}
