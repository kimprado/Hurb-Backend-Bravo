package config

import (
	"flag"
	"time"

	"github.com/jinzhu/configor"
)

// config possui configuracões da aplicação
var config = Configuration{}
var loaded bool

// Configuration - type
type Configuration struct {
	Environment struct {
		Name string `default:"dev"`
	}

	Server struct {
		Port string `default:"8081"`
	}

	RedisDB RedisDB

	Logging struct {
		File  string `required:"false"`
		Level LoggingLevels
	}

	CurrencyManager struct {
		SupportedCurrencies []string `required:"true"`
	}

	RatesFinder struct {
		BaseCurrency string        `required:"true"`
		EntryTimeout time.Duration `required:"true"`
	}
}

// Redis representa configuração de conexão Redis
type Redis struct {
	Host     string `default:"localhost"`
	Port     int    `default:"6379"`
	User     string `required:"false"`
	Password string `required:"false"`
	// Prefixo de todas chaves
	Prefix string `default:"exchange"`
}

// RedisDB representa configuração Redis em modo DB
type RedisDB Redis

// NewRedisDB cria novo RedisDB
func NewRedisDB(c Configuration) (r RedisDB) {
	r = c.RedisDB
	return
}

// LoggingLevels representa loggers e seus respectivos níveis
type LoggingLevels map[string]string

// NewLoggingLevels cria novo LoggingLevels
func NewLoggingLevels(c Configuration) (ll LoggingLevels) {
	ll = c.Logging.Level
	return
}

// NewConfig -
func NewConfig(configLocationFile string) (c Configuration, err error) {
	if !loaded {
		var configLocation string
		if configLocationFile != "" {
			configLocation = configLocationFile
		} else {
			configLocation = loadFlags()
		}
		config, err = loadConfig(configLocation)
		if err != nil {
			return
		}
		c = config
		config = c
		loaded = true
		return
	}
	c = config
	return
}

func loadFlags() (configPath string) {
	cp := flag.String("config-location", "./configs/config-dev.json", "Caminho para arquivo de configuração")

	flag.Parse()

	configPath = *cp
	return
}

func loadConfig(configLocation string) (config Configuration, err error) {
	configApp := new(Configuration)

	cfg := configor.New(&configor.Config{
		ENVPrefix: "EXCHANGE",
	})

	err = cfg.Load(configApp, configLocation)

	if err != nil {
		return
	}
	config = *configApp
	return
}
