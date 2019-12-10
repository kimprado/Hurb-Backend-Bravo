package api

import (
	"fmt"
	"net/http"

	"github.com/rep/exchange/internal/pkg/commom/errors"
	"github.com/rep/exchange/internal/pkg/currencyexchange"
)

func statusCode(e error) (s int) {

	switch v := e.(type) {
	case *currencyexchange.LookupCurrencyError:
		s = http.StatusServiceUnavailable
	case *currencyexchange.UnsupportedCurrencyError:
		s = http.StatusBadRequest
	case *errors.ParametersError:
		s = http.StatusBadRequest
	default:
		panic(fmt.Sprintf("Tipo de erro não definido %T", v))
	}
	return
}
