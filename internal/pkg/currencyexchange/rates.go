package currencyexchange

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
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

func (q Quote) float64() (f float64) {
	f = (float64)(q)
	return
}

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
	url := urlQuery(rf.baseCurrency, cs...)
	err = sendRequest(url, &dto)

	if err == errHTTPStatusBadRequest {
		pmError := newRateQuoteExternalServiceParametersError(url)
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
	locker  sync.RWMutex
	service RatesFinder
	cache   map[string]*rateChacheEntry
	timeout Timeout
	logger  logging.LoggerRatesCache
}

// NewRatesFinderCache é responsável por instanciar RatesFinderCache
func NewRatesFinderCache(s *RatesFinderService, cfg config.Configuration, l logging.LoggerRatesCache) (rf *RatesFinderCache) {
	rf = new(RatesFinderCache)
	rf.service = s
	rf.logger = l
	rf.timeout = (Timeout)(cfg.RatesFinder.EntryTimeout * time.Second)
	rf.cache = make(map[string]*rateChacheEntry)
	return
}

// Find delega para outras implementações. Consulta Taxas de Câmbio para moedas
func (rf *RatesFinderCache) Find(cs ...Currency) (r map[string]*Rate, err error) {

	rf.rlock()
	rates := make(map[string]*Rate)
	mustFind := []Currency{}

	for _, c := range cs {
		if v, ok := rf.cache[c.Code()]; ok && v.Active() {
			rates[c.Code()] = v.Rate()
			continue
		}
		mustFind = append(mustFind, c)
	}
	if len(cs) == len(rates) {
		r = rates
		rf.logger.Tracef("Cotações recuperadas do cache %q\n", r)
		defer rf.runlock()
		return
	}

	rf.runlock()

	newerRates, err := rf.service.Find(mustFind...)

	rf.lock()
	defer rf.unlock()

	for _, r := range newerRates {
		if v, ok := rf.cache[r.currency.Code()]; ok && !v.Active() {
			v.renew(r, rf.timeout)
		} else {
			rf.cache[r.currency.Code()] = newRateChacheEntry(rf, r, rf.timeout, rf.logger)
		}
		rates[r.currency.Code()] = r
	}
	r = rates

	return
}

func (rf *RatesFinderCache) lock() {
	rf.locker.Lock()
}

func (rf *RatesFinderCache) unlock() {
	rf.locker.Unlock()
}

func (rf *RatesFinderCache) rlock() {
	rf.locker.RLock()
}

func (rf *RatesFinderCache) runlock() {
	rf.locker.RUnlock()
}

// Timeout indica duração da entrada no chache.
// Usado na invalidação da entrada no cache.
type Timeout time.Duration

func (t Timeout) duration() (d time.Duration) {
	d = (time.Duration)(t)
	return
}

// rateChacheEntry representa entrada no chache para Rate
type rateChacheEntry struct {
	cache   *RatesFinderCache
	rate    *Rate
	active  bool
	timeout Timeout
	logger  logging.LoggerRatesCache
}

// newRateChacheEntry cria instância de rateChacheEntry
func newRateChacheEntry(rf *RatesFinderCache, r *Rate, t Timeout, l logging.LoggerRatesCache) (entry *rateChacheEntry) {
	entry = new(rateChacheEntry)
	entry.cache = rf
	entry.rate = r
	entry.active = true
	entry.timeout = t
	entry.logger = l
	entry.verifyAndMakeInvalid()

	if entry.logger.IsInfoEnabled() {
		until := time.Now().Add(entry.timeout.duration()).Format("2006-01-02 15:04:05")
		entry.logger.Infof("Nova Entrada %q válida até %q por %q\n", entry.rate.currency.Code(), until, entry.timeout.duration())
	}
	return
}

func (entry *rateChacheEntry) Rate() (r *Rate) {
	r = entry.rate
	return
}

func (entry *rateChacheEntry) Active() (a bool) {
	a = entry.active
	return
}

func (entry *rateChacheEntry) renew(r *Rate, t Timeout) {
	entry.rate = r
	entry.active = true
	entry.timeout = t

	if entry.logger.IsTraceEnabled() {
		until := time.Now().Add(entry.timeout.duration()).Format("2006-01-02 15:04:05")
		entry.logger.Tracef("Entrada renovada %q válida até %q por %q\n", entry.rate.currency.Code(), until, entry.timeout.duration())
	}

	entry.verifyAndMakeInvalid()
	return
}

func (entry *rateChacheEntry) verifyAndMakeInvalid() {
	go func() {
		time.Sleep(entry.timeout.duration())
		entry.cache.lock()
		defer entry.cache.unlock()
		entry.active = false
		entry.logger.Tracef("Entrada %q removida do cache\n", entry.rate.currency.Code())
	}()
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
func (rf *RatesFinderProxy) Find(cs ...Currency) (rates map[string]*Rate, err error) {
	rates, err = rf.service.Find(cs...)
	return
}
