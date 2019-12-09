package currencyexchange

import "github.com/google/wire"

// PkgSet define providers do pacote
var PkgSet = wire.NewSet(
	NewCalculatorController,
	// Define que a implementação Padão de Calculator é CalculatorController
	wire.Bind(new(Calculator), new(*CalculatorController)),
	NewCurrencyManagerProxy,
	// Define que a implementação Padão de CurrencyManager é CurrencyManagerProxy
	wire.Bind(new(CurrencyManager), new(*CurrencyManagerProxy)),
)
