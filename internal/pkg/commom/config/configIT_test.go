// +build test integration

package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

var now = time.Now()

const configTemplate = `
{
	"environment": {
		"name": "test-%s"
	},
    "server": {
        "port": "3080"
    },
    "redisDB": {
        "host": "host-IT-test",
        "port": 6379
	},
    "currencymanager": {
        "SUPPORTEDCURRENCIES": [
            "USD"
            ,"BRL"
            ,"EUR"
            ,"BTC"
            ,"ETH"
        ]
    },
    "ratesfinder": {
        "basecurrency": "USD",
        "entrytimeout": 60
    },
    "logging": {
        "level": {
            "ROOT": "INFO"
        }
    }
}
`

func TestLoadConfig(t *testing.T) {
	f, d, err := createTmpFile()

	if err != nil {
		t.Fatalf("Falha ao criar arquivo temporário para teste %v\n", err)
	}
	defer cleanResources(f, d)

	dateTime := now.Format("2006-01-02 15:04:05")
	writeFile(f, fmt.Sprintf(configTemplate, dateTime))

	expect := struct {
		environment       string
		serverPort        string
		redisDbHost       string
		redisDbPort       int
		basecurrency      string
		defaultcurrencies []string
		entrytimeout      time.Duration
		logging           map[string]string
	}{
		environment:       "test-" + dateTime,
		serverPort:        "3080",
		redisDbHost:       "host-IT-test",
		redisDbPort:       6379,
		basecurrency:      "USD",
		defaultcurrencies: []string{"USD", "BRL", "EUR", "BTC", "ETH"},
		entrytimeout:      60,
		logging: map[string]string{
			"ROOT": "INFO",
		},
	}

	var c Configuration

	c, err = loadConfig(f.Name())

	if err != nil {
		t.Errorf("Erro ao carregar configurações %v", err)
		return
	}

	if expect.environment != c.Environment.Name {
		t.Errorf("Environment esperado %q é diferente de %q\n", expect.environment, c.Environment.Name)
	}
	if expect.serverPort != c.Server.Port {
		t.Errorf("serverPort esperado %q é diferente de %q\n", expect.serverPort, c.Server.Port)
	}
	if expect.redisDbHost != c.RedisDB.Host {
		t.Errorf("redisDbHost esperado %q é diferente de %q\n", expect.redisDbHost, c.RedisDB.Host)
	}
	if expect.redisDbPort != c.RedisDB.Port {
		t.Errorf("redisDbPort esperado %q é diferente de %q\n", expect.redisDbPort, c.RedisDB.Port)
	}
	if expect.basecurrency != c.RatesFinder.BaseCurrency {
		t.Errorf("basecurrency esperado %q é diferente de %q\n", expect.basecurrency, c.RatesFinder.BaseCurrency)
	}
	if fmt.Sprintf("%v", expect.defaultcurrencies) != fmt.Sprintf("%v", c.CurrencyManager.SupportedCurrencies) {
		t.Errorf("defaultcurrencies esperado %q é diferente de %q\n", expect.defaultcurrencies, c.CurrencyManager.SupportedCurrencies)
	}
	if expect.entrytimeout != c.RatesFinder.EntryTimeout {
		t.Errorf("entrytimeout esperado %q é diferente de %q\n", expect.entrytimeout, c.RatesFinder.EntryTimeout)
	}

	for k, v := range expect.logging {
		z, ok := c.Logging.Level[k]
		if !ok {
			t.Errorf("Log de nível %q não encontrado na lista\n", k)
		}
		if ok && v != z {
			t.Errorf("Log esperado %q[%s] é diferente de %q[%s]\n", k, v, k, z)
		}
	}

	if err := f.Close(); err != nil {
		t.Fatalf("Falha ao fechar arquivo temporário para teste %v\n", err)
	}
}

func cleanResources(tmpFile *os.File, tempDir string) {
	os.Remove(tmpFile.Name())
	os.Remove(tempDir)
}

func createTmpFile() (tmpFile *os.File, tempDir string, err error) {

	tempDir, err = ioutil.TempDir("", "challenge-exchange-api")
	if err != nil {
		return
	}
	tmpFile, err = ioutil.TempFile(tempDir, "config-testing-*.json")
	return

}

func writeFile(tmpFile *os.File, content string) (err error) {

	text := []byte(content)
	if _, err = tmpFile.Write(text); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}
	return
}
