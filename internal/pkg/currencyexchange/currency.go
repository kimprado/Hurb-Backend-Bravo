package currencyexchange

// CurrencyManager gerencia moedas disponíveis para cálculo.
type CurrencyManager interface {
	CurrencyAdder
	CurrencyRemover
	Find() Currency
}

// Currency representa entidade monetária com valor associado
type Currency struct {
	code string
	rate float64
}
