// +build test unit

package currencyexchange

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExchangeAmount(t *testing.T) {
	t.Parallel()

	var err error

	currencyManager := &CurrencyManagerMock{func(currency string) (c *Currency, err error) {
		return &Currency{currency}, nil
	}}

	ratesFinder := &RatesFinderMock{func(currencies ...Currency) (rates map[string]*Rate, err error) {
		rates = make(map[string]*Rate)
		for _, c := range currencies {
			var quote Quote
			switch c.Code() {
			case "BRL":
				quote = Quote(4.2179556517)
				break
			case "EUR":
				quote = Quote(0.9013881377)
				break
			case "CAD":
				quote = Quote(1.3237923251)
				break
			case "USD":
				quote = Quote(1)
				break
			default:
				return nil, errors.New("Moeda não configurada no teste")
			}
			rates[c.Code()] = &Rate{c, quote}
		}
		return
	}}

	calculator, err := initializeCalculatorControllerTest(currencyManager, ratesFinder)
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
		{1.068513057192417, "BRL", "EUR", 5},
		{1.5692345230875773, "BRL", "CAD", 5},
		{1.1854083856914917, "BRL", "USD", 5},
		{46.79400000162659, "EUR", "BRL", 10},
		{31.862668877320868, "CAD", "BRL", 10},
		{42.179556516999995, "USD", "BRL", 10},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			t.Parallel()

			result, err := calculator.Exchange(tc.curFrom, tc.curTo, tc.amount)
			if err != nil {
				t.Errorf("Erro inesperado ao realizar Exchange do TestCase %+v: %v", tc, err)
				return
			}
			assert.NotNil(t, result)
			assert.Equal(t, tc.expect, result.value)
		})
	}

}

type CurrencyManagerMock struct {
	f func(currency string) (c *Currency, err error)
}

func (cm *CurrencyManagerMock) Find(currency string) (c *Currency, err error) {
	return cm.f(currency)
}
func (cm *CurrencyManagerMock) Add(CurrencyDTO) (err error) {
	return
}
func (cm *CurrencyManagerMock) Remove(CurrencyDTO) (err error) {
	return
}

type RatesFinderMock struct {
	f func(currencies ...Currency) (rates map[string]*Rate, err error)
}

func (rf *RatesFinderMock) Find(currencies ...Currency) (rates map[string]*Rate, err error) {
	return rf.f(currencies...)
}
