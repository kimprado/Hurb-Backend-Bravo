// +build test integration

package currencyexchange

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindSupportedCurrency(t *testing.T) {
	t.Parallel()

	var err error

	err = setUpSupportedCurrency(t)
	if err != nil {
		t.Errorf("Erro ao preparar teste: %+v\n", err)
		return
	}

	c, err := initializeConfigTest()
	if err != nil {
		t.Errorf("Erro ao criar Configuração: %+v\n", err)
		return
	}

	c.RedisDB.Prefix = t.Name()

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
	assert.Equal(t, result.Code(), "BRL")

	result, err = currencyManager.Find("BEF")
	if err != nil {
		t.Errorf("Consulta Redis %v\n", err)
		return
	}

	assert.Nil(t, result)

}

func TestAddSupportedCurrency(t *testing.T) {
	t.Parallel()

	var err error

	c, err := initializeConfigTest()
	if err != nil {
		t.Errorf("Erro ao criar Configuração: %+v\n", err)
		return
	}

	c.RedisDB.Prefix = t.Name()

	currencyManager, err := initializeCurrencyManagerDBTest(c)
	if err != nil {
		t.Errorf("Conexão banco de dados %v\n", err)
		return
	}
	assert.NotNil(t, currencyManager)

	err = currencyManager.Add(CurrencyDTO{"BRL"})
	if err != nil {
		t.Errorf("Consulta Redis %v\n", err)
		return
	}

}

func TestRemoveSupportedCurrency(t *testing.T) {
	t.Parallel()

	var err error

	err = setUpSupportedCurrency(t)
	if err != nil {
		t.Errorf("Erro ao preparar teste: %+v\n", err)
		return
	}

	c, err := initializeConfigTest()
	if err != nil {
		t.Errorf("Erro ao criar Configuração: %+v\n", err)
		return
	}

	c.RedisDB.Prefix = t.Name()

	currencyManager, err := initializeCurrencyManagerDBTest(c)
	if err != nil {
		t.Errorf("Conexão banco de dados %v\n", err)
		return
	}
	assert.NotNil(t, currencyManager)

	err = currencyManager.Remove(CurrencyDTO{"BRL"})
	if err != nil {
		t.Errorf("Consulta Redis %v\n", err)
		return
	}

}

// setUpSupportedCurrency cria carga de dados para teste.
// Popula Redis com valores em nova chave para o teste.
// Nome da chave se baseia no nome do teste, o que permite
// executar testes de integração em paralelo :).
func setUpSupportedCurrency(t *testing.T) (err error) {
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

	_, err = con.Do("SADD", t.Name()+":currency:supported", currencyBRL)

	if err != nil {
		return
	}

	return
}
