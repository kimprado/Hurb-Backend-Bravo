package api

import (
	"fmt"
	"net/http"

	"github.com/rep/exchange/internal/pkg/commom/errors"
	"github.com/rep/exchange/internal/pkg/currencyexchange"
)

func statusCode(e error) (s int) {

	switch v := e.(type) {
	case *currencyexchange.CurrencyCreationParametersError:
		s = http.StatusBadRequest // 400
	case *currencyexchange.CurrencyCreationError:
		s = http.StatusInternalServerError // 500
	case *currencyexchange.RemoveCurrencyNotFoundError:
		s = http.StatusNotFound // 404
	case *currencyexchange.RemoveCurrencyError:
		s = http.StatusInternalServerError // 500
	case *currencyexchange.LookupCurrencyError:
		s = http.StatusServiceUnavailable // 503
	case *currencyexchange.UnsupportedCurrencyError:
		s = http.StatusBadRequest // 400
	case *currencyexchange.RateQuoteExternalServiceParametersError:
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
		// Adotar alguma das seguintes medidas:
		// 	- Aplicar testes/verificações.
		// 	- Eventualmente substituir por log.err, caso não exista.
		// 	- cobertura suficiente que garanta que não será executado.
		panic(fmt.Sprintf("Tipo de erro não definido %T", v))
	}
	return
}
