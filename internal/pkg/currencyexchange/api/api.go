package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rep/exchange/internal/pkg/commom/errors"
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
	var err error

	queryValues := req.URL.Query()

	paramErr := errors.NewParametersError()

	from := queryValues.Get("from")
	to := queryValues.Get("to")
	strAmount := queryValues.Get("amount")

	if from == "" {
		paramErr.Add(
			errors.ParameterError{
				Name:   "from",
				Value:  from,
				Reason: "'from' não pode ser vazio",
			},
		)
	}
	if to == "" {
		paramErr.Add(
			errors.ParameterError{
				Name:   "to",
				Value:  to,
				Reason: "'to' não pode ser vazio",
			},
		)
	}
	if strAmount == "" {
		paramErr.Add(
			errors.ParameterError{
				Name:   "amount",
				Value:  strAmount,
				Reason: "'amount' não pode ser vazio",
			},
		)
	}

	amount, err := strconv.ParseFloat(strAmount, 64)
	if strAmount != "" && err != nil {
		paramErr.Add(
			errors.ParameterError{
				Name:   "amount",
				Value:  strAmount,
				Reason: fmt.Sprintf("Não foi possivel converter %q. Informe valor monetário. Ex: 123.45", strAmount),
			},
		)
	}

	if paramErr.ContainsError() {
		v.logger.Warnf("Consulta Exchange com parâmetro(s) inválido(s):\n\t%v\n", queryValues)

		web.NewHTTPResponse(
			res,
			http.StatusBadRequest,
			nil,
			paramErr,
		).WriteJSON()

		return
	}

	err = v.calculator.Exchange(from, to, amount)

	if err != nil {
		v.logger.Warnf("Erro ao realizar Exchange: %v\n", err)

		web.NewHTTPResponse(
			res,
			http.StatusBadRequest,
			nil,
			err,
		).WriteJSON()

		return
	}

	web.NewHTTPResponse(
		res,
		http.StatusOK,
		struct{}{},
		nil,
	).WriteJSON()

}
