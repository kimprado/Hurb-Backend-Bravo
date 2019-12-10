package currencyexchange

// RatesDTO mapeia resposta da API de câmbio
type RatesDTO struct {
	Rates map[string]Quote `json:"rates"`
}
