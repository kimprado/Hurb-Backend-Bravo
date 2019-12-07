package logging

import (
	l "github.com/kimprado/sllog/pkg/logging"
	"github.com/rep/exchange/internal/pkg/commom/config"
)

//Logger para logar ROOT
type Logger struct {
	l.Logger
}

//LoggerRedisDB para logar infra.redis.db
type LoggerRedisDB struct {
	l.Logger
}

//LoggerWebServer para logar webserver
type LoggerWebServer struct {
	l.Logger
}

//NewLogger cria Logger ""(ROOT)
func NewLogger(configLevels config.LoggingLevels) (log Logger) {
	log = Logger{l.NewLogger("", configLevels)}
	return
}

//NewRedisDB cria Logger "infra.redis.db"
func NewRedisDB(configLevels config.LoggingLevels) (log LoggerRedisDB) {
	log = LoggerRedisDB{l.NewLogger("infra.redis.db", configLevels)}
	return
}

//NewWebServer cria Logger "webserver"
func NewWebServer(configLevels config.LoggingLevels) (log LoggerWebServer) {
	log = LoggerWebServer{l.NewLogger("webserver", configLevels)}
	return
}
