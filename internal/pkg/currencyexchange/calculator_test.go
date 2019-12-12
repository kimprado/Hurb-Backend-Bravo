// +build test unit

package currencyexchange

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculatorControllerExchangeAmount(t *testing.T) {
	t.Parallel()

	var err error

	currencyManager := &CurrencyManagerMock{func(currency string) (c *Currency, err error) {
		return &Currency{currency}, nil
	}}

	calculator, err := initializeCalculatorControllerTest(currencyManager, nil)
	if err != nil {
		t.Errorf("Criação serviço %v\n", err)
		return
	}
	assert.NotNil(t, calculator)

	testCases := []struct {
		expect  float64
		curFrom string
		curTo   string
		amount  float64
	}{
		{5, "BRL", "EUR", 1.068513057},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			t.Parallel()

			result, err := calculator.Exchange(tc.curFrom, tc.curTo, tc.amount)

			if err != nil {
				t.Errorf("Erro inesperado ao realizar Exchange de %#v - %v", tc, err)
				return
			}

			assert.NotNil(t, result)
			assert.Equal(t, tc.expect, result)

		})
	}

}

type CurrencyManagerMock struct {
	f func(currency string) (c *Currency, err error)
}

func (cm *CurrencyManagerMock) Find(currency string) (c *Currency, err error) {
	return cm.f(currency)
}
func (cm *CurrencyManagerMock) Add(currency string) (err error) {
	return
}
func (cm *CurrencyManagerMock) Remove(currency string) (err error) {
	return
}
