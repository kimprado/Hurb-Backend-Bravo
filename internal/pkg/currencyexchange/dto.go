package currencyexchange

// CurrencyDTO representa moeda troca de informações da API
type CurrencyDTO struct {
	Code string `json:"code"`
}

// RatesDTO mapeia resposta da API de câmbio
type RatesDTO struct {
	Rates map[string]Quote `json:"rates"`
}

// ExchangeValueDTO representa DTO de ExchangeValue
type ExchangeValueDTO struct {
	Value float64 `json:"value"`
}

// NewExchangeValueDTO Cria instância de ExchangeValueDTO
func NewExchangeValueDTO(ex *ExchangeValue) (dto *ExchangeValueDTO) {
	dto = new(ExchangeValueDTO)
	dto.Value = ex.Value()
	return
}
