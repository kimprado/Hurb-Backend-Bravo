package currencyexchange

// CurrencyManager gerencia moedas disponíveis para cálculo.
type CurrencyManager interface {
	CurrencyAdder
	CurrencyRemover
}
