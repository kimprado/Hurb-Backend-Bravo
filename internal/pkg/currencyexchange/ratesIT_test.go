// +build test integration

package currencyexchange

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindRatesQuote(t *testing.T) {
	t.Parallel()

	var err error

	c, err := initializeConfigTest()
	if err != nil {
		t.Errorf("Erro ao criar Configuração: %+v\n", err)
		return
	}

	ratesFinder, err := initializeRatesFinderServiceTest(c)
	if err != nil {
		t.Errorf("Criação serviço %v\n", err)
		return
	}
	assert.NotNil(t, ratesFinder)

	testCases := []struct {
		curFrom     string
		curTo       string
		errExpected error
	}{
		{"BRL", "EUR", nil},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			result, err := ratesFinder.Find(Currency{tc.curFrom}, Currency{tc.curTo})

			if tc.errExpected == nil && err != nil {
				t.Errorf("Error inesperado %v", err)
				return
			}

			if got := errors.Is(err, tc.errExpected); !got && tc.errExpected != nil {
				t.Errorf("Esperado erro %v, mas obtido %v", tc.errExpected, err)
				return
			}

			var r *Rate
			r, _ = result[tc.curFrom]
			assert.NotNil(t, r)
			assert.Equal(t, tc.curFrom, r.currency.Code())
			assert.NotZero(t, r.quote)

			r, _ = result[tc.curTo]
			assert.NotNil(t, r)
			assert.Equal(t, tc.curTo, r.currency.Code())
			assert.NotZero(t, r.quote)

		})
	}

}
