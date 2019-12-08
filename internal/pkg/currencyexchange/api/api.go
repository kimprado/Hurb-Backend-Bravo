package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rep/exchange/internal/pkg/commom/logging"
	"github.com/rep/exchange/internal/pkg/currencyexchange"
)

// Controller trata requisições http de paredão
type Controller struct {
	calculator currencyexchange.Calculator
	logger     logging.LoggerAPIExchange
}

// NewController é responsável por instanciar Controller
func NewController(l logging.LoggerAPIExchange) (r *Controller) {
	r = new(Controller)
	r.logger = l
	return
}

// Exchange calcula câmbio
func (v *Controller) Exchange(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	v.calculator.Exchange()

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(http.StatusOK)

	json.NewEncoder(res).Encode(
		struct {
			Value float64
		}{
			Value: 123.45,
		},
	)

}
