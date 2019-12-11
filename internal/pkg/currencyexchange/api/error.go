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
		s = http.StatusServiceUnavailable // 503
	case *currencyexchange.UnsupportedCurrencyError:
		s = http.StatusBadRequest // 400
	case *currencyexchange.RateQuoteServiceParametersError:
		s = http.StatusBadRequest // 400
	case *currencyexchange.CurrencyRateQuoteNotFoundError:
		s = http.StatusServiceUnavailable // 503
	case *currencyexchange.RateQuoteNotFoundError:
		s = http.StatusInternalServerError // 500
	case *currencyexchange.LookupRatesQuoteError:
		s = http.StatusServiceUnavailable // 503
	case *errors.ParametersError:
		s = http.StatusBadRequest // 400
	default:
		// TODO: Fazer com que este código não seja executado.
		// Aplicar testes/verificações.
		// Eventualmente substituir por log.err.
		panic(fmt.Sprintf("Tipo de erro não definido %T", v))
	}
	return
}
