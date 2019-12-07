package redis

import (
	"github.com/google/wire"
	"github.com/rep/exchange/internal/pkg/commom/config"
	"github.com/rep/exchange/internal/pkg/commom/logging"
)

var pkgSetConfigTest = wire.NewSet(
	newIntegrationConfigFile,
	config.PkgSet,
)

var pkgSetTest = wire.NewSet(
	PkgSet,
	config.NewLoggingLevels,
	config.NewRedisDB,
	logging.PkgSet,
)

func newIntegrationConfigFile() string {
	return "../../../../configs/config-integration.json"
}
