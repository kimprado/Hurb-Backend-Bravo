package currencyexchange

import (
	"fmt"

	"github.com/rep/exchange/internal/pkg/commom/errors"
)

// LookupCurrencyError representa erro na pesquisa de moeda
type LookupCurrencyError struct {
	*errors.GenericError
}

// newLookupCurrencyError cria instância de LookupCurrencyError
func newLookupCurrencyError(currency string) (e *LookupCurrencyError) {
	e = new(LookupCurrencyError)
	e.GenericError = errors.NewGenericError(
		"Falha na consulta de moeda",
		fmt.Sprintf("Não foi possível consultar moeda %q", currency),
	)
	return
}

func (e *LookupCurrencyError) Error() string {
	return e.GenericError.Error()
}

// UnsupportedCurrencyError representa erro na pesquisa de moeda
type UnsupportedCurrencyError struct {
	*errors.ParametersError
}

// newUnsupportedCurrencyError cria instância de UnsupportedCurrencyError
func newUnsupportedCurrencyError(param, value, currency string) (e *UnsupportedCurrencyError) {
	e = new(UnsupportedCurrencyError)

	e.ParametersError = errors.NewParametersError()
	e.ParametersError.Title = "Uma ou Mais moedas não são válidas para câmbio"
	e.ParametersError.Add(
		errors.ParameterError{
			Name:   param,
			Value:  value,
			Reason: fmt.Sprintf("Moeda %q indisponível para câmbio", currency),
		},
	)
	return
}

func (e *UnsupportedCurrencyError) Error() string {
	return e.ParametersError.Error()
}

// RateQuoteNotFoundError representa erro na consulta de cotação
type RateQuoteNotFoundError struct {
	*errors.ParametersError
}

// newUnsupportedCurrencyError cria instância de RateQuoteNotFoundError
func newRateQuoteNotFoundError() (e *RateQuoteNotFoundError) {
	e = new(RateQuoteNotFoundError)
	e.ParametersError = errors.NewParametersError()
	e.ParametersError.Title = "Falha na consulta consulta de Uma ou Mais Quotaçõs"
	return
}

// AddQuote inclui erro de cotacação
func (e *RateQuoteNotFoundError) AddQuote(currency string) string {
	e.ParametersError.Add(
		errors.ParameterError{
			Name:   "",
			Value:  currency,
			Reason: fmt.Sprintf("Cotação da moeda %q indisponível", currency),
		},
	)
	return e.ParametersError.Error()
}

func (e *RateQuoteNotFoundError) Error() string {
	return e.ParametersError.Error()
}
