package currencyexchange

// Exchanger realiza cálculo de conversão monetária.
type Exchanger interface {
	// Calcula câmbio monetário
	Exchange(from, to Rate, amount float64) (v *ExchangeValue, err error)
}

// ExchangeValue representa valor convertido
type ExchangeValue struct {
	value float64
}

// Value retorna valor de value
func (c ExchangeValue) Value() (cd float64) {
	cd = c.value
	return
}

// CurrencyExchanger implementa Exchanger.
// Realiza conversão monetária com taxas de câmbio.
type CurrencyExchanger struct {
}

// NewCurrencyExchanger cria instância de CurrencyExchanger
func NewCurrencyExchanger() (c *CurrencyExchanger) {
	c = new(CurrencyExchanger)
	return
}

// Exchange converte valor na moeda destino usando cotações fornecidas
func (c *CurrencyExchanger) Exchange(from, to Rate, amount float64) (v *ExchangeValue, err error) {
	ex := amount / from.quote.float64() * to.quote.float64()
	v = &ExchangeValue{ex}
	return
}
