// +build test integration

package currencyexchange

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindSupportedCurrency(t *testing.T) {
	t.Parallel()

	var err error

	setUpFindSupportedCurrency()
	if err != nil {
		t.Errorf("Erro ao preparar teste: %+v\n", err)
		return
	}

	c, err := initializeConfigTest()
	if err != nil {
		t.Errorf("Erro ao criar Configuração: %+v\n", err)
		return
	}

	currencyManager, err := initializeCurrencyManagerDBTest(c)
	if err != nil {
		t.Errorf("Conexão banco de dados %v\n", err)
		return
	}
	assert.NotNil(t, currencyManager)

	result, err := currencyManager.Find("BRL")
	if err != nil {
		t.Errorf("Consulta Redis %v\n", err)
		return
	}

	assert.NotNil(t, result)

}

func setUpFindSupportedCurrency() (err error) {
	currencyBRL := "BRL"

	c, err := initializeConfigTest()
	if err != nil {
		return
	}

	redis, err := initializeRedisTest(c)
	if err != nil {
		return
	}

	con := redis.Get()

	_, err = con.Do("SADD", "currency:supported", currencyBRL)

	if err != nil {
		return
	}

	return
}
