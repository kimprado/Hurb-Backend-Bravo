package currencyexchange

import "github.com/rep/exchange/internal/pkg/commom/logging"

// CurrencyManager gerencia moedas disponíveis para cálculo.
type CurrencyManager interface {
	CurrencyAdder
	CurrencyRemover

	// Consulta moedas ativas
	Find(currency string) (c *Currency)
}

// Currency representa entidade monetária com valor associado
type Currency struct {
	code string
	rate float64
}

// CurrencyManagerProxy implementa proxy para CurrencyManagers
type CurrencyManagerProxy struct {
	logger logging.LoggerCurrency
}

// NewCurrencyManagerProxy é responsável por instanciar Controller
func NewCurrencyManagerProxy(l logging.LoggerCurrency) (c *CurrencyManagerProxy) {
	c = new(CurrencyManagerProxy)
	c.logger = l
	return
}

// Find delega para outras implementações. Consulta moedas ativas
func (cm *CurrencyManagerProxy) Find(currency string) (c *Currency) {

	return
}

// Add delega para outras implementações. Adiciona moeda
func (cm *CurrencyManagerProxy) Add(currency string) {

}

// Remove delega para outras implementações. Remove moeda
func (cm *CurrencyManagerProxy) Remove(currency string) {

}
