package currencyexchange

import (
	"fmt"
	"regexp"

	"github.com/rep/exchange/internal/pkg/commom/config"
	"github.com/rep/exchange/internal/pkg/commom/logging"
	"github.com/rep/exchange/internal/pkg/infra/redis"
)

var rgUpperCase = regexp.MustCompile(`^[A-Z]{1,}$`)

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
}

func newCurrency(code string) (c *Currency) {
	c = new(Currency)
	c.code = code
	return
}

func newCurrencyFromDTO(dto CurrencyDTO) (c *Currency) {
	c = new(Currency)
	c.code = dto.Code
	return
}

// Code retorna valor de code
func (c Currency) Code() (cd string) {
	cd = c.code
	return
}

// Valid retorna erros que invalidam criação da moeda
// Retorna err == nil caso não exista inconsistência
func (c Currency) Valid() (err error) {

	paramErr := newCurrencyCreationParametersError()
	if len(c.code) != 3 {
		paramErr.Add(
			"code",
			c.code,
			fmt.Sprintf("Código da moeda %q inválido. Código deve ter 3 caracteres. Ex: USD", c.code),
		)
	}
	if !rgUpperCase.MatchString(c.code) {
		paramErr.Add(
			"code",
			c.code,
			fmt.Sprintf("Código da moeda %q inválido. Código deve ter apenas letras maiúsculas. Ex: USD", c.code),
		)
	}
	if paramErr.ContainsError() {
		err = paramErr
	}

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
func (cm *CurrencyManagerProxy) Add(dto CurrencyDTO) (err error) {
	err = cm.db.Add(dto)
	return
}

// Remove delega para outras implementações. Remove moeda
func (cm *CurrencyManagerProxy) Remove(dto CurrencyDTO) (err error) {
	err = cm.db.Remove(dto)
	return
}

// CurrencyManagerDB implementa CurrencyManager com acesso a DB
type CurrencyManagerDB struct {
	redisClient redis.DBConnection
	redisCfg    config.RedisDB
	cfg         config.Configuration
	logger      logging.LoggerCurrency
}

// NewCurrencyManagerDB é responsável por instanciar Controller
func NewCurrencyManagerDB(r redis.DBConnection, cr config.RedisDB, cfg config.Configuration, l logging.LoggerCurrency) (c *CurrencyManagerDB) {
	c = new(CurrencyManagerDB)
	c.redisClient = r
	c.redisCfg = cr
	c.cfg = cfg
	c.logger = l
	return
}

// Find delega para outras implementações. Consulta moedas ativas
func (cm *CurrencyManagerDB) Find(currency string) (c *Currency, err error) {

	//TODO: Criar repositório para retirar código de infraestrutura da camada de serviço
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

// Add inclui moeda
func (cm *CurrencyManagerDB) Add(dto CurrencyDTO) (err error) {

	c := newCurrencyFromDTO(dto)
	err = c.Valid()
	if err != nil {
		return
	}

	var currency string
	currency = c.Code()

	//TODO: Criar repositório para retirar código de infraestrutura da camada de serviço
	con := cm.redisClient.Get()
	defer con.Close()

	_, err = con.Do("SADD", fmt.Sprintf("%s:currency:supported", cm.redisCfg.Prefix), currency)

	if err != nil {
		return
	}

	return
}

// Remove deleta moeda
func (cm *CurrencyManagerDB) Remove(dto CurrencyDTO) (err error) {

	var currency string
	currency = dto.Code

	//TODO: Criar repositório para retirar código de infraestrutura da camada de serviço
	const notFound = 0

	con := cm.redisClient.Get()
	defer con.Close()

	var reply interface{}
	reply, err = con.Do("SREM", fmt.Sprintf("%s:currency:supported", cm.redisCfg.Prefix), currency)

	if err != nil {
		return
	}

	if reply.(int64) == notFound {
		err = newRemoveCurrencyNotFoundError(currency)
	}

	return
}

// LoadDefaultSupportedCurrencies carrega moedas suportadas por padrão.
// Realiza carga apenas 1 vez.
func (cm *CurrencyManagerDB) LoadDefaultSupportedCurrencies() (err error) {

	//TODO: Criar repositório para retirar código de infraestrutura da camada de serviço
	con := cm.redisClient.Get()
	defer con.Close()

	var loaded bool

	loaded, err = cm.verifyLoaded()
	if err != nil {
		return
	}
	if loaded {
		return
	}

	err = con.Send("MULTI")
	if err != nil {
		return
	}

	con.Send("SET", fmt.Sprintf("%s:config:currency:loaded", cm.redisCfg.Prefix), 1)
	for _, c := range cm.cfg.CurrencyManager.SupportedCurrencies {
		con.Send("SADD", fmt.Sprintf("%s:currency:supported", cm.redisCfg.Prefix), c)
	}

	_, err = con.Do("EXEC")
	if err != nil {
		return
	}

	return
}

func (cm *CurrencyManagerDB) verifyLoaded() (loaded bool, err error) {

	//TODO: Criar repositório para retirar código de infraestrutura da camada de serviço
	const found = 1

	con := cm.redisClient.Get()
	defer con.Close()

	var reply interface{}
	reply, err = con.Do("EXISTS", fmt.Sprintf("%s:config:currency:loaded", cm.redisCfg.Prefix))

	if err != nil {
		return
	}

	if reply.(int64) == found {
		loaded = true
	}

	return
}
