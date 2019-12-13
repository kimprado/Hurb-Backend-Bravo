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
		{"BRL", "USD", nil},
		{"EUR", "BRL", nil},
		{"EUR", "USD", nil},
		{"USD", "BRL", nil},
		{"USD", "EUR", nil},
		{"", "BRL", &RateQuoteExternalServiceParametersError{}},
		{"BRL", "", &RateQuoteExternalServiceParametersError{}},
		{"BRLs", "EUR", &RateQuoteExternalServiceParametersError{}},
		{"EUR", "BRLs", &RateQuoteExternalServiceParametersError{}},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			t.Parallel()

			result, err := ratesFinder.Find(Currency{tc.curFrom}, Currency{tc.curTo})

			if err != nil && tc.errExpected == nil {
				t.Errorf("Erro inesperado %v", err)
				return
			}

			got := errors.Is(err, tc.errExpected)

			if got && tc.errExpected != nil {
				return
			}

			if !got && tc.errExpected != nil {
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
