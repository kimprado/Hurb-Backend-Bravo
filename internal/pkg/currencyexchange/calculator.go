package currencyexchange

import "github.com/rep/exchange/internal/pkg/commom/logging"

// Calculator é um ponto de entrada para comportamento
// da aplicação. O controlador do domínio.
type Calculator interface {
	Exchange()
}

// CurrencyAdder é um ponto de entrada para comportamento
// da aplicação. O controlador do domínio.
type CurrencyAdder interface {
	Add()
}

// CurrencyRemover é um ponto de entrada para comportamento
// da aplicação. O controlador do domínio.
type CurrencyRemover interface {
	Remove()
}

// CalculatorController é o ponto de entrada para comportamentos
// da aplicação. Implementação do controlador do domínio.
type CalculatorController struct {
	logger logging.LoggerCalculator
}

// NewCalculatorController é responsável por instanciar Controller
func NewCalculatorController(l logging.LoggerCalculator) (c *CalculatorController) {
	c = new(CalculatorController)
	c.logger = l
	return
}

// Exchange executa conversão monetária
func (c *CalculatorController) Exchange() {

	c.logger.Debugf("Calculando câmbio")

}
