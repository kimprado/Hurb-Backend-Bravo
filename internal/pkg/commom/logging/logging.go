package logging

import (
	l "github.com/kimprado/sllog/pkg/logging"
	"github.com/rep/exchange/internal/pkg/commom/config"
)

// Logger para logar ROOT
type Logger struct {
	l.Logger
}

// LoggerAPIExchange para logar currencyexchange.api.exchange
type LoggerAPIExchange struct {
	l.Logger
}

// LoggerCalculator para logar currencyexchange.calculator
type LoggerCalculator struct {
	l.Logger
}

// LoggerCurrency para logar currencyexchange.currency
type LoggerCurrency struct {
	l.Logger
}

// LoggerRates para logar currencyexchange.rates
type LoggerRates struct {
	l.Logger
}

// LoggerRedisDB para logar infra.redis.db
type LoggerRedisDB struct {
	l.Logger
}

// LoggerWebServer para logar webserver
type LoggerWebServer struct {
	l.Logger
}

// NewLogger cria Logger ""(ROOT)
func NewLogger(configLevels config.LoggingLevels) (log Logger) {
	log = Logger{l.NewLogger("", configLevels)}
	return
}

// NewLoggerAPIExchange cria Logger "currencyexchange.api.exchange"
func NewLoggerAPIExchange(configLevels config.LoggingLevels) (log LoggerAPIExchange) {
	log = LoggerAPIExchange{l.NewLogger("currencyexchange.api.exchange", configLevels)}
	return
}

// NewCalculator cria Logger "currencyexchange.calculator"
func NewCalculator(configLevels config.LoggingLevels) (log LoggerCalculator) {
	log = LoggerCalculator{l.NewLogger("currencyexchange.calculator", configLevels)}
	return
}

// NewCurrency cria Logger "currencyexchange.currency"
func NewCurrency(configLevels config.LoggingLevels) (log LoggerCurrency) {
	log = LoggerCurrency{l.NewLogger("currencyexchange.currency", configLevels)}
	return
}

// NewRates cria Logger "currencyexchange.rates"
func NewRates(configLevels config.LoggingLevels) (log LoggerRates) {
	log = LoggerRates{l.NewLogger("currencyexchange.currency", configLevels)}
	return
}

// NewRedisDB cria Logger "infra.redis.db"
func NewRedisDB(configLevels config.LoggingLevels) (log LoggerRedisDB) {
	log = LoggerRedisDB{l.NewLogger("infra.redis.db", configLevels)}
	return
}

// NewWebServer cria Logger "webserver"
func NewWebServer(configLevels config.LoggingLevels) (log LoggerWebServer) {
	log = LoggerWebServer{l.NewLogger("webserver", configLevels)}
	return
}
