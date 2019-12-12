package app

import (
	"github.com/rep/exchange/internal/pkg/commom/logging"
	"github.com/rep/exchange/internal/pkg/currencyexchange"
	"github.com/rep/exchange/internal/pkg/webserver"
)

// ExchangeApp representa instância da aplicação
type ExchangeApp struct {
	cm        *currencyexchange.CurrencyManagerDB
	webServer *webserver.WebServer
	logger    logging.Logger
}

// NewExchangeApp cria app
func NewExchangeApp(cm *currencyexchange.CurrencyManagerDB, ws *webserver.WebServer, l logging.Logger) (a *ExchangeApp) {
	a = new(ExchangeApp)
	a.webServer = ws
	a.cm = cm
	a.logger = l
	return
}

// Bootstrap é responsável por iniciar a aplicação
func (a *ExchangeApp) Bootstrap() {
	a.logger.Infof("Iniciando serviços da aplicação...\n")

	a.cm.LoadDefaultSupportedCurrencies()

	a.webServer.Start()
}
