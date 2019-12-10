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
