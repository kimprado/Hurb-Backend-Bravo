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

// CurrencyRateQuoteNotFoundError representa erro de cotação não encontrada para moeda
type CurrencyRateQuoteNotFoundError struct {
	*errors.ParametersError
}

// newCurrencyRateQuoteNotFoundError cria instância de CurrencyRateQuoteNotFoundError
func newCurrencyRateQuoteNotFoundError() (e *CurrencyRateQuoteNotFoundError) {
	e = new(CurrencyRateQuoteNotFoundError)
	e.ParametersError = errors.NewParametersError()
	e.ParametersError.Title = "Falha na consulta de Uma ou Mais Cotações"
	return
}

// Add inclui erro de cotacação
func (e *CurrencyRateQuoteNotFoundError) Add(currency string) {
	e.ParametersError.Add(
		errors.ParameterError{
			Name:   "",
			Value:  currency,
			Reason: fmt.Sprintf("Cotação da moeda %q indisponível", currency),
		},
	)
}

func (e *CurrencyRateQuoteNotFoundError) Error() string {
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

// Is informa se target == e. Verifica se e é do tipo
// RateQuoteServiceParametersError, DomainError.
func (e *RateQuoteServiceParametersError) Is(target error) bool {
	switch target.(type) {
	case *RateQuoteServiceParametersError:
		return true
	case *errors.DomainError:
		return true
	default:
		return false
	}
}

// LookupRatesQuoteError representa erro na pesquisa de Taxas de Câmbio
type LookupRatesQuoteError struct {
	*errors.GenericError
}

// newLookupRatesQuoteError cria instância de LookupRatesQuoteError
func newLookupRatesQuoteError() (e *LookupRatesQuoteError) {
	e = new(LookupRatesQuoteError)
	e.GenericError = errors.NewGenericError(
		"Falha na consulta de Taxas de Câmbio",
		fmt.Sprintf("Não foi possível consultar Taxas de Câmbio"),
	)
	return
}

func (e *LookupRatesQuoteError) Error() string {
	return e.GenericError.Error()
}
