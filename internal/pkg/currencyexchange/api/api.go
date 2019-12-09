package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rep/exchange/internal/pkg/commom/logging"
	"github.com/rep/exchange/internal/pkg/commom/web"
	"github.com/rep/exchange/internal/pkg/currencyexchange"
)

// Controller trata requisições http de paredão
type Controller struct {
	calculator currencyexchange.Calculator
	logger     logging.LoggerAPIExchange
}

// NewController é responsável por instanciar Controller
func NewController(c currencyexchange.Calculator, l logging.LoggerAPIExchange) (r *Controller) {
	r = new(Controller)
	r.calculator = c
	r.logger = l
	return
}

// Exchange calcula câmbio
func (v *Controller) Exchange(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	queryValues := req.URL.Query()

	from := queryValues.Get("from")
	to := queryValues.Get("to")
	strAmount := queryValues.Get("amount")

	amount, err := strconv.ParseFloat(strAmount, 64)

	if err != nil {
		v.logger.Errorf("Consulta Exchange %+v\n", params)
		res.Header().Set("Content-Type", "application/json; charset=UTF-8")

		web.NewHTTPResponse(
			res,
			http.StatusBadRequest,
			nil,
			web.NewErrorMessage(fmt.Sprintf("Amount %q is invalid", strAmount)),
		).WriteJSON()

		return
	}

	v.calculator.Exchange(from, to, amount)

	value := struct {
		Value float64
	}{
		Value: 123.45,
	}

	hr := web.NewHTTPResponse(res, http.StatusOK, value, nil)
	hr.WriteJSON()

}
