// +build test unit

package currencyexchange

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculatorControllerExchangeAmount(t *testing.T) {
	t.Parallel()

	var err error

	//TODO: Praparar Mock Objects dos Serviços
	calculator, err := initializeCalculatorControllerTest(nil, nil)
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
