package currencyexchange

import (
	"github.com/rep/exchange/internal/pkg/commom/logging"
)

// Calculator é um ponto de entrada para comportamento
// da aplicação. O controlador do domínio.
type Calculator interface {
	// Calcula câmbio monetário
	Exchange(from, to string, amount float64) (err error)
}

// CurrencyAdder é um ponto de entrada para comportamento
// da aplicação. O controlador do domínio.
type CurrencyAdder interface {
	Add(currency string) (err error)
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
	logger logging.LoggerCalculator
}

// NewCalculatorController é responsável por instanciar Controller
func NewCalculatorController(cm CurrencyManager, l logging.LoggerCalculator) (c *CalculatorController) {
	c = new(CalculatorController)
	c.cm = cm
	c.logger = l
	return
}

// Exchange executa conversão monetária
func (c *CalculatorController) Exchange(from, to string, amount float64) (err error) {
	// Validar moedas disponíveis para conversão
	// Consultar taxas de câmbio
	// Calcular conversão

	curFrom, err := c.cm.Find(from)
	if err != nil {
		err = newLookupCurrencyError(from)
		return
	}

	c.logger.Debugf("Calculando câmbio com moeda: %v\n", curFrom)

	return
}
