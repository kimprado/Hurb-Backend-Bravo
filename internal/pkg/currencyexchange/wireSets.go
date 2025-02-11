package currencyexchange

import (
	"github.com/google/wire"
	"github.com/rep/exchange/internal/pkg/commom/config"
	"github.com/rep/exchange/internal/pkg/commom/logging"
	"github.com/rep/exchange/internal/pkg/infra/redis"
)

// PkgSet define providers do pacote
var PkgSet = wire.NewSet(
	NewCalculatorController,
	// Define que a implementação Padão de Calculator é CalculatorController
	wire.Bind(new(Calculator), new(*CalculatorController)),
	// Define que a implementação Padão de CurrencyAdder é CalculatorController
	wire.Bind(new(CurrencyAdder), new(*CalculatorController)),
	// Define que a implementação Padão de CurrencyRemover é CalculatorController
	wire.Bind(new(CurrencyRemover), new(*CalculatorController)),
	NewCurrencyManagerProxy,
	// Define que a implementação Padão de CurrencyManager é CurrencyManagerProxy
	wire.Bind(new(CurrencyManager), new(*CurrencyManagerProxy)),
	NewCurrencyManagerDB,
	NewRatesFinderProxy,
	// Define que a implementação Padão de RatesFinder é RatesFinderProxy
	wire.Bind(new(RatesFinder), new(*RatesFinderProxy)),
	NewRatesFinderCache,
	NewRatesFinderService,
	NewCurrencyExchanger,
	// Define que a implementação Padão de Exchanger é CurrencyExchanger
	wire.Bind(new(Exchanger), new(*CurrencyExchanger)),
)

var pkgSetConfigTest = wire.NewSet(
	newIntegrationConfigFile,
	config.PkgSet,
)

var pkgSetTest = wire.NewSet(
	PkgSet,
	config.NewLoggingLevels,
	config.NewRedisDB,
	logging.PkgSet,
	redis.PkgSet,
)

var pkgSetCalculatorMocksTest = wire.NewSet(
	NewCalculatorController,
	NewCurrencyExchanger,
	// Define que a implementação Padão de Exchanger é CurrencyExchanger
	wire.Bind(new(Exchanger), new(*CurrencyExchanger)),
	pkgSetConfigTest,
	logging.PkgSet,
)

func newIntegrationConfigFile() string {
	return "../../../configs/config-integration.json"
}
