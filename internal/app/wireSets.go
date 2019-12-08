package app

import (
	"github.com/google/wire"
	"github.com/rep/exchange/internal/pkg/commom/config"
	"github.com/rep/exchange/internal/pkg/commom/logging"
	"github.com/rep/exchange/internal/pkg/currencyexchange/api"
	"github.com/rep/exchange/internal/pkg/infra/redis"
	"github.com/rep/exchange/internal/pkg/webserver"
)

//AppSet define providers do pacote
var AppSet = wire.NewSet(
	config.PkgSet,
	logging.PkgSet,
	api.PkgSet,
	redis.PkgSet,
	webserver.PkgSet,
	NewExchangeApp,
)
