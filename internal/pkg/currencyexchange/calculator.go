package currencyexchange

import (
	"github.com/rep/exchange/internal/pkg/commom/errors"
	"github.com/rep/exchange/internal/pkg/commom/logging"
)

// Calculator é um ponto de entrada para comportamento
// da aplicação. O controlador do domínio.
type Calculator interface {
	// Calcula câmbio monetário
	Exchange(from, to string, amount float64) (ex *ExchangeValue, err error)
}

// CurrencyAdder é um ponto de entrada para comportamento
// da aplicação. O controlador do domínio.
type CurrencyAdder interface {
	Add(dto CurrencyDTO) (err error)
}

// CurrencyRemover é um ponto de entrada para comportamento
// da aplicação. O controlador do domínio.
type CurrencyRemover interface {
	Remove(currency string) (err error)
}

// CalculatorController é o ponto de entrada para comportamentos
// da aplicação. Implementação do controlador do domínio.
type CalculatorController struct {
	cm     CurrencyManager
	rf     RatesFinder
	ex     Exchanger
	logger logging.LoggerCalculator
}

// NewCalculatorController é responsável por instanciar Controller
func NewCalculatorController(cm CurrencyManager, rf RatesFinder, ex Exchanger, l logging.LoggerCalculator) (c *CalculatorController) {
	c = new(CalculatorController)
	c.cm = cm
	c.rf = rf
	c.ex = ex
	c.logger = l
	return
}

// Exchange executa conversão monetária
func (c *CalculatorController) Exchange(from, to string, amount float64) (ex *ExchangeValue, err error) {
	// Calcular conversão

	curFrom, err := c.cm.Find(from)
	if err != nil {
		err = newLookupCurrencyError(from)
		return
	}
	if curFrom == nil {
		err = newUnsupportedCurrencyError("from", from, from)
		return
	}

	curTo, err := c.cm.Find(to)
	if err != nil {
		err = newLookupCurrencyError(to)
		return
	}
	if curTo == nil {
		err = newUnsupportedCurrencyError("to", to, to)
		return
	}

	c.logger.Tracef("Calculando câmbio com moedas: %v | %v\n", curFrom, curTo)

	rates, err := c.rf.Find(*curFrom, *curTo)
	if err != nil {
		err = errors.GetDomainErrorOr(err, newLookupRatesQuoteError())
		return
	}
	if len(rates) != 2 {
		err = newRateQuoteNotFoundError()
		return
	}

	rFrom := rates[from]
	rTo := rates[to]

	c.logger.Tracef("Taxas de câmbio: %v\n", rates)

	ex, err = c.ex.Exchange(*rFrom, *rTo, amount)

	return
}

// Add executa inclusão de moeda suportada
func (c *CalculatorController) Add(dto CurrencyDTO) (err error) {
	err = c.cm.Add(dto)
	if err != nil {
		err = errors.GetDomainErrorOr(err, newSupportedCurrencyCreationError())
	}
	return
}
