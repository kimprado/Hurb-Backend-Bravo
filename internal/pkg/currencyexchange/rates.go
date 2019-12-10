package currencyexchange

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rep/exchange/internal/pkg/commom/config"
	"github.com/rep/exchange/internal/pkg/commom/logging"
)

const baseURL = "https://api.exchangeratesapi.io/latest?base=%s&symbols=%s"

// RatesFinder pesquisa taxas de câmbio.
type RatesFinder interface {
	Find(currencies ...Currency) (rates map[string]*Rate, err error)
}

// Quote representa valor de uma cotação
type Quote float64

// BaseCurrency representa moeda de lastro
type BaseCurrency string

// Rate mapeia resposta da API de câmbio
type Rate struct {
	currency Currency
	quote    Quote
}

// NewRate cria instância de Rate
func NewRate(c Currency, q Quote) (r *Rate) {
	r = new(Rate)
	r.currency = c
	r.quote = q
	return
}

// RatesFinderService responsável por consultar cotações de câmbio
type RatesFinderService struct {
	baseCurrency BaseCurrency
	logger       logging.LoggerCurrency
}

// NewRatesFinderService é responsável por criar RatesFinderService
func NewRatesFinderService(c config.Configuration, l logging.LoggerCurrency) (rf *RatesFinderService) {
	rf = new(RatesFinderService)
	rf.baseCurrency = (BaseCurrency)(c.RatesFinder.BaseCurrency)
	rf.logger = l
	return
}

// Find pesquisa cotação online
func (rf *RatesFinderService) Find(cs ...Currency) (rates map[string]*Rate, err error) {
	rates = make(map[string]*Rate)

	var dto RatesDTO
	err = sendRequest(urlQuery(rf.baseCurrency, cs...), &dto)

	if err != nil {
		return
	}

	nfError := newRateQuoteNotFoundError()
	for _, c := range cs {

		quote, ok := dto.Rates[c.Code()]
		if !ok {
			nfError.AddQuote(c.Code())
		}
		rates[c.Code()] = NewRate(c, quote)
	}

	if nfError.ContainsError() {
		err = nfError
	}

	return
}

func sendRequest(url string, result interface{}) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	return err
}

func urlQuery(b BaseCurrency, cs ...Currency) (u string) {
	currencyCodes := make([]string, 0, len(cs))
	for _, c := range cs {
		currencyCodes = append(currencyCodes, c.Code())
	}
	u = fmt.Sprintf(baseURL, b, strings.Join(currencyCodes, ","))
	return
}
