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

// newRateQuoteNotFoundError cria instância de RateQuoteNotFoundError
func newRateQuoteNotFoundError() (e *RateQuoteNotFoundError) {
	e = new(RateQuoteNotFoundError)
	e.ParametersError = errors.NewParametersError()
	e.ParametersError.Title = "Falha na consulta de Uma ou Mais Cotações"
	return
}

// Add inclui erro de cotacação
func (e *RateQuoteNotFoundError) Add(currency string) {
	e.ParametersError.Add(
		errors.ParameterError{
			Name:   "",
			Value:  currency,
			Reason: fmt.Sprintf("Cotação da moeda %q indisponível", currency),
		},
	)
}

func (e *RateQuoteNotFoundError) Error() string {
	return e.ParametersError.Error()
}

// RateQuoteServiceParametersError representa erro de parâmetros na consulta
// ao serviço de cotação
type RateQuoteServiceParametersError struct {
	*errors.ParametersError
}

// newRateQuoteServiceParametersError cria instância de RateQuoteServiceParametersError
func newRateQuoteServiceParametersError() (e *RateQuoteServiceParametersError) {
	e = new(RateQuoteServiceParametersError)
	e.ParametersError = errors.NewParametersError()
	e.ParametersError.Title = "Um ou Mais parâmetros não são válidos no Serviço Externo de Cotações"
	return
}

// Add inclui erro de cotacação
func (e *RateQuoteServiceParametersError) Add(currency string) {
	e.ParametersError.Add(
		errors.ParameterError{
			Name:   "",
			Value:  currency,
			Reason: fmt.Sprintf("Verificar se moeda %q é válida para cotação ou lastro", currency),
		},
	)
}

func (e *RateQuoteServiceParametersError) Error() string {
	return e.ParametersError.Error()
}

// Is informa se target == e
func (e *RateQuoteServiceParametersError) Is(target error) bool {
	_, ok := target.(*RateQuoteServiceParametersError)
	if !ok {
		return false
	}
	return true
}
