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
