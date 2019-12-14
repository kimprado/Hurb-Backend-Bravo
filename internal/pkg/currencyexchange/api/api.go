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
	calculator      currencyexchange.Calculator
	currencyAdder   currencyexchange.CurrencyAdder
	currencyRemover currencyexchange.CurrencyRemover
	logger          logging.LoggerAPIExchange
}

// NewController é responsável por instanciar Controller
func NewController(c currencyexchange.Calculator, a currencyexchange.CurrencyAdder, r currencyexchange.CurrencyRemover, l logging.LoggerAPIExchange) (ctrl *Controller) {
	ctrl = new(Controller)
	ctrl.calculator = c
	ctrl.currencyAdder = a
	ctrl.currencyRemover = r
	ctrl.logger = l
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
		v.logger.Warnf("Consulta Exchange : %v\n", paramErr)

		web.NewHTTPResponse(
			res,
			statusCode(paramErr),
			nil,
			paramErr,
		).WriteJSON()

		return
	}

	var ex *currencyexchange.ExchangeValue
	ex, err = v.calculator.Exchange(from, to, amount)

	if err != nil {
		v.logger.Warnf("Erro ao realizar Exchange: %v\n", err)

		web.NewHTTPResponse(
			res,
			statusCode(err),
			nil,
			err,
		).WriteJSON()

		return
	}

	dto := currencyexchange.NewExchangeValueDTO(ex)

	web.NewHTTPResponse(
		res,
		http.StatusOK,
		dto,
		nil,
	).WriteJSON()

}

// AddSupportedCurrency cria nova moeda suportada
func (v *Controller) AddSupportedCurrency(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	var err error
	var dto currencyexchange.CurrencyDTO

	currency := params.ByName("currency")
	dto = currencyexchange.CurrencyDTO{
		Code: currency,
	}

	err = v.currencyAdder.Add(dto)

	if err != nil {
		v.logger.Errorf("Erro ao criar moeda: %v\n", err)

		web.NewHTTPResponse(
			res,
			statusCode(err),
			nil,
			err,
		).WriteJSON()

		return
	}

	web.NewHTTPResponse(
		res,
		http.StatusNoContent,
		nil,
		nil,
	).WriteJSON()

}

// RemoveSupportedCurrency remove moeda suportada
func (v *Controller) RemoveSupportedCurrency(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	var err error
	var dto currencyexchange.CurrencyDTO

	currency := params.ByName("currency")
	dto = currencyexchange.CurrencyDTO{
		Code: currency,
	}

	err = v.currencyRemover.Remove(dto)

	if err != nil {
		v.logger.Warnf("Erro ao remover moeda: %v\n", err)

		web.NewHTTPResponse(
			res,
			statusCode(err),
			nil,
			err,
		).WriteJSON()

		return
	}

	web.NewHTTPResponse(
		res,
		http.StatusNoContent,
		nil,
		nil,
	).WriteJSON()

}
