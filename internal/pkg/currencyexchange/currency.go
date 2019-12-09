package currencyexchange

import (
	"fmt"

	"github.com/rep/exchange/internal/pkg/commom/config"
	"github.com/rep/exchange/internal/pkg/commom/logging"
	"github.com/rep/exchange/internal/pkg/infra/redis"
)

// CurrencyManager gerencia moedas disponíveis para cálculo.
type CurrencyManager interface {
	CurrencyAdder
	CurrencyRemover

	// Consulta moedas ativas
	Find(currency string) (c *Currency, err error)
}

// Currency representa entidade monetária com valor associado
type Currency struct {
	code string
	rate float64
}

func newCurrency(code string) (c *Currency) {
	c = new(Currency)
	c.code = code
	return
}

// Code retorna valor de code
func (c Currency) Code() (cd string) {
	cd = c.code
	return
}

// CurrencyManagerProxy implementa proxy para CurrencyManagers
type CurrencyManagerProxy struct {
	db     CurrencyManager
	logger logging.LoggerCurrency
}

// NewCurrencyManagerProxy é responsável por instanciar Controller
func NewCurrencyManagerProxy(db *CurrencyManagerDB, l logging.LoggerCurrency) (c *CurrencyManagerProxy) {
	c = new(CurrencyManagerProxy)
	c.db = db
	c.logger = l
	return
}

// Find delega para outras implementações. Consulta moedas ativas
func (cm *CurrencyManagerProxy) Find(currency string) (c *Currency, err error) {
	c, err = cm.db.Find(currency)
	return
}

// Add delega para outras implementações. Adiciona moeda
func (cm *CurrencyManagerProxy) Add(currency string) (err error) {

	return
}

// Remove delega para outras implementações. Remove moeda
func (cm *CurrencyManagerProxy) Remove(currency string) {

}

// CurrencyManagerDB implementa CurrencyManager com acesso a DB
type CurrencyManagerDB struct {
	redisClient redis.DBConnection
	redisCfg    config.RedisDB
	logger      logging.LoggerCurrency
}

// NewCurrencyManagerDB é responsável por instanciar Controller
func NewCurrencyManagerDB(r redis.DBConnection, cr config.RedisDB, l logging.LoggerCurrency) (c *CurrencyManagerDB) {
	c = new(CurrencyManagerDB)
	c.redisClient = r
	c.redisCfg = cr
	c.logger = l
	return
}

// Find delega para outras implementações. Consulta moedas ativas
func (cm *CurrencyManagerDB) Find(currency string) (c *Currency, err error) {
	const found = 1

	con := cm.redisClient.Get()
	defer con.Close()

	var reply interface{}
	reply, err = con.Do("SISMEMBER", fmt.Sprintf("%s:currency:supported", cm.redisCfg.Prefix), currency)

	if err != nil {
		return
	}

	if reply.(int64) == found {
		c = newCurrency(currency)
	}

	return
}

// Add delega para outras implementações. Adiciona moeda
func (cm *CurrencyManagerDB) Add(currency string) (err error) {

	con := cm.redisClient.Get()
	defer con.Close()

	_, err = con.Do("SADD", fmt.Sprintf("%s:currency:supported", cm.redisCfg.Prefix), currency)

	if err != nil {
		return
	}

	return
}

// Remove delega para outras implementações. Remove moeda
func (cm *CurrencyManagerDB) Remove(currency string) {

}
