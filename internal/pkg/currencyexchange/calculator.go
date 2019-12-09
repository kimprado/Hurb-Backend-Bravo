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
	Add(currency string)
}

// CurrencyRemover é um ponto de entrada para comportamento
// da aplicação. O controlador do domínio.
type CurrencyRemover interface {
	Remove(currency string)
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
func (c *CalculatorController) Exchange() {
	// Validar moedas disponíveis para conversão
	// Consultar taxas de câmbio
	// Calcular conversão

	cur1, _ := c.cm.Find("")
	c.logger.Debugf("Calculando câmbio com moeda: %v\n", cur1)
}
