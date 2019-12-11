package currencyexchange

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rep/exchange/internal/pkg/commom/config"
	"github.com/rep/exchange/internal/pkg/commom/logging"
)

const baseURL = "https://api.exchangeratesapi.io/latest?base=%s&symbols=%s"

var errHTTPStatusBadRequest = errors.New("HTTP Bad Request")

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

func (r *Rate) String() (s string) {
	s = fmt.Sprintf("%v %.6f", r.currency.Code(), (float64)(r.quote))
	return
}

// RatesFinderService responsável por consultar cotações de câmbio
// Usa serviço externo (exchangeratesapi.io)
type RatesFinderService struct {
	baseCurrency BaseCurrency
	logger       logging.LoggerRates
}

// NewRatesFinderService é responsável por criar RatesFinderService
func NewRatesFinderService(c config.Configuration, l logging.LoggerRates) (rf *RatesFinderService) {
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

	if err == errHTTPStatusBadRequest {
		pmError := newRateQuoteServiceParametersError()
		pmError.Add((string)(rf.baseCurrency))
		for _, c := range cs {
			pmError.Add(c.Code())
		}
		err = pmError
		return
	}
	if err != nil {
		return
	}

	nfError := newCurrencyRateQuoteNotFoundError()
	for _, c := range cs {

		quote, ok := dto.Rates[c.Code()]
		if !ok {
			nfError.Add(c.Code())
		}
		rates[c.Code()] = NewRate(c, quote)
	}

	if nfError.ContainsError() {
		err = nfError
	}

	return
}

func sendRequest(url string, result interface{}) error {
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusBadRequest {
		return errHTTPStatusBadRequest
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

// RatesFinderCache implementa proxy para RatesFinders.
// Realiza cache das consultas.
type RatesFinderCache struct {
	service RatesFinder
	logger  logging.LoggerRates
}

// NewRatesFinderCache é responsável por instanciar RatesFinderCache
func NewRatesFinderCache(s *RatesFinderService, l logging.LoggerRates) (c *RatesFinderCache) {
	c = new(RatesFinderCache)
	c.service = s
	c.logger = l
	return
}

// Find delega para outras implementações. Consulta Taxas de Câmbio para moedas
func (cm *RatesFinderCache) Find(cs ...Currency) (rates map[string]*Rate, err error) {
	rates, err = cm.service.Find(cs...)
	return
}

// RatesFinderProxy implementa proxy para RatesFinders
type RatesFinderProxy struct {
	service RatesFinder
	logger  logging.LoggerRates
}

// NewRatesFinderProxy é responsável por instanciar RatesFinderProxy
func NewRatesFinderProxy(s *RatesFinderCache, l logging.LoggerRates) (c *RatesFinderProxy) {
	c = new(RatesFinderProxy)
	c.service = s
	c.logger = l
	return
}

// Find delega para outras implementações. Consulta Taxas de Câmbio para moedas
func (cm *RatesFinderProxy) Find(cs ...Currency) (rates map[string]*Rate, err error) {
	rates, err = cm.service.Find(cs...)
	return
}
