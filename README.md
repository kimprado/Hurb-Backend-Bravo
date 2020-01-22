# Documentação Currency Exchange API

Descrição da solução para o desafio da conversão monetária([Desafio Bravo](https://github.com/hurbcom/challenge-bravo))

- [O Problema](#O-Problema)
    - [Back-end](#Back-end)
        - [Dependências](#Dependências)
            - [Boilerplate Code](#Boilerplate-Code)
    - [Documentação API](#Documentação-API)
- [Instalação e Execução](#Instalação-e-Execução)
- [Ambiente Desenvolvimento](#Ambiente-Desenvolvimento)
    - [Primeira Execução](#Primeira-Execução)
        - [Instalação das Dependências](#Instalação-das-Dependências)
    - [Execução](#Execução)
    - [Infra Desenvolvimento](#Infra-Desenvolvimento)
    - [Infra Testes](#Infra-Testes)
    - [Infra Documentação](#Infra-Documentação)
- [Testes](#Testes)
    - [Unitários](#Unitários)
    - [Integração](#Integração)
    - [Unitários e Integração](#tests-all)
    - [Carga](#Carga)
- [Empacotamento](#Empacotamento)
- [Comandos Make](#comandos-make)
- [Melhorias](#Melhorias)

## O Problema

Foi implementada solução que permite fazer conversão monetária de valores entre diferentes moedas, tendo moeda padrão de lastro USD(parametrizável).
Também é possível adicionar e remover dinamicamente moedas suportadas, por meio de chamadas a um endpoint específico

A seguir temos exemplos de utilização da API. Para uma documentação mais completa verifique o tópico [swagger](swagger).

 - Conversão de Valores
 
 Ex: Converter: 5,00 BRL em EUR

 ```sh
 curl -X GET 'http://localhost:3000/exchange?from=BRL&to=EUR&amount=5'
 ```

Teremos a resposta JSON que o valor **BRL 5** coresponde a **EUR 1.094954449941698**:

 ```json
{
    "value": 1.094954449941698
}
 ```

 - Criação de moeda suportada

 Ex: Permitir suporte a moeda CAD

 ```sh
 curl -X PUT 'http://localhost:3000/currencies/CAD'
 ```

 - Remover suporte a moeda

 Ex: Remover suporte a moeda CAD

 ```sh
 curl -X DELETE 'http://localhost:3000/currencies/CAD'
 ```

### Back-end

Implementação de API como microserviço.

Aplicação tem uma camada de API que faz a interface HTTP, sendo um dos endpoint REST.
Foi definida uma camada de negócio no pacote *currencyexchange* que tem 
um ponto de entrada, *CalculatorController*. 

Uma chamada HTTP típica para câmbio monetário segue este caminho, iniciando pelo [*Controller.Exchange()*](internal/pkg/currencyexchange/api/api.go) da API.

- *CalculatorController* do negócio, colabora com *CurrencyManagerDB* para acessar moedas ativas armazenadas no Redis.

- Na sequência *CalculatorController* consulta cotações reais em serviço externo com *RatesFinderService*, que tem suas respostas gardadas em cache em memória, por um Proxy.

- Por último a conversão é realizada por meio de *CurrencyExchanger*.

*CalculatorController* implementa *Calculator* e depende de interfaces na colaboração com outros serviços.

A seguir outra forma de representar Interações entre as interfaçes que os componentes implementam. *Controller(api)* é o único participante concreto.

```
Controller(api)
  \--(1)--Exchange()---> Calculator
                           \--(2)--Find()-------> CurrencyManager
                           \--(3)--Find()-------> RatesFinder
                           \--(4)--Exchange()---> Exchanger
```

Segue descrição dos pacotes e arquivos da solução.

 ```sh
 tree -L 5
.
├── cmd
│   └── exchangeAPI
│       ├── main.go             # Main da Aplicação                           
│       ├── wire_gen.go         # _boilerplate code_
│       └── wire.go             # Provedores de Dependência da aplicão
├── configs
│   ├── config.env              # Arquivo ENV de configuração usado em 'make run'
├── go.mod                      # Dependências da aplicação
├── internal
│   ├── app
│   │   ├── exchangeAPI.go      # Inicialização da Aplicação
│   │   └── wireSets.go         # Declaração das depências DI
│   └── pkg
│       ├── commom              # Package de utilitários
│       ├── currencyexchange    # Package principal do contexto de negócio
│       │   ├── api             # Package da API HTTP da aplicação
│       │   │   ├── api.go
│       │   │   ├── error.go    # Mapeia erros de negócio para erros HTTP Status Code
│       │   │   └── wireSets.go 
│       │   ├── calculator.go   # Ponto de entrada do negócio. Atende solicitações da API
│       │   ├── calculator_test.go
│       │   ├── currency.go     # Implementa entidade Moeda
│       │   ├── currencyIT_test.go
│       │   ├── dto.go          # Implementa DTO de usado para recebimento de dados
│       │   ├── errors.go       # Erros de negócio
│       │   ├── rates.go        # Implementa entidade Taxas de Câmbio
│       │   ├── ratesIT_test.go
│       │   ├── rule.go         # Implementa regras de convesão monetária
│       │   ├── wire_gen.go     # _boilerplate code_
│       │   ├── wire.go         # Provedores de Dependência dos testes da package
│       │   └── wireSets.go     # Declaração das depências DI de testes 
│       ├── infra               # Pacote de infraestrutura
│       └── webserver           
│           ├── home.go
│           ├── webserver.go    # Implementa servidor HTTP que expões API de negócio
│           └── wireSets.go
 ```

#### Dependências

- **[Go](http://golang.org/)** - Liguagem usada na implementação da API. 
- **[Wire](http://github.com/google/wire)** - por que é mais prático usar injeção de dependências.
    - Com Wire aplicação não quebra em runtime por falha na declaração de dependências.
    - Ciclo de desenvolvimento é mais rápido por não precisar iniciar aplicação para testar dependências.
    - Mensagens de falha na resolução do grafo de dependências são claras.
    - Ponto desfavorável: Gera arquivos Boilerplate Code ***wire_gen.go*** que devem ser commitados
- **[Redigo](http://github.com/gomodule/redigo)** - driver performárico, que mantém API do Redis.
- **[slLog](http://github.com/kimprado/sllog)** - escrevi esta lib para configurar logging como no Spring Boot.

##### Boilerplate Code

A dependência Wire gera arquivos ***wire_gen.go***.

### Documentação API

Implementada documentação interativa com [swagger](api/swagger.yml). Para acessá-la execute o ambiente como descrito a seguir([Instalação e Execução](#Instalação-e-Execução)), e depois siga as instruções em [Infra Documentação](#Infra-Documentação).


## Instalação e Execução

Para fazer deploy e execução do projeto rode os seguintes comandos:

```sh
./configure
make run
```

Ao final na execução o comando printa no console as urls para acesso aos serviço.

- http://localhost:80/      (nginx)   - Página web com links úteis
- http://0.0.0.0:3000/      (API)     - URL da API
- http://localhost:80/docs  (swagger) - DOcumentação interativa

Ex:

```sh
make run
...
Acesse nginx:
http://localhost:80/
Acesse API:
http://0.0.0.0:3000/
Acesse swagger:
http://localhost:80/docs
...
```

## Ambiente Desenvolvimento

Segue como instalar e configurar o ambiente e ferramentas de desenvolvimento do projeto.

### Primeira Execução

#### Instalação das Dependências

Execute script *[configure](configure)* presente na raiz do repositório para fazer download e instalação das dependências. O script também cria, caso necessário, a configuração da IDE e de execução da aplicação.

```sh
./configure
```

Após executar feche e abra outro terminal.

As seguintes ferramentas serão provisionadas:

- **Docker** - ferramenta usada para containerização.
- **Docker Compose** - ferramenta usada para orquestração em ambiende de dev.
- **Go** - linguagem de programação.

Instale o Wire com o seguinte comando:

```sh
go get github.com/google/wire/cmd/wire@v0.3.0
```

Os seguintes arquivos são criados pelo arquivo configure:

- .vscode/settings.json - Arquivo da IDE VSCode
- .vscode/launch.json - Arquivo da IDE VSCode
- configs/config.json - Configuração opcional da aplicação em tempo de desenvolvimento
- configs/config.env - Configurações padão de desenvolvimento e  injetadas como Env Vars no deploy do Docker Compose.

### Execução

 - Executar solução

```sh
make run
```

 - Interromper a execução

```sh
make stop
```

### Infra Desenvolvimento

- Iniciar infra de Desenvolvimento

```sh
make infra-start
```

- Interromper infra de Desenvolvimento
```sh
make infra-stop
```

### Infra Testes

- Iniciar infra de Testes

```sh
docker-compose up -d --build redis-test
```

- Interromper infra de Testes
```sh
docker-compose rm -fsv redis-test
```

### Infra Documentação

- http://localhost:80/docs - Modo interativo com back-end
    ```sh
    make run
    ```

- http://localhost:8080/ - Somente documentação
    ```sh
    docker-compose up -d swagger
    ```

- http://localhost:80/docs - Modo interativo com back-end rodando pela IDE, por exemplo
    ```sh
    make infra-start
    ```

## Testes

Fizemos separação dos testes em dois grupos principais, *[Unitários](#Unitários)* e *[Integração](#Integração)*. Os grupos de testes são separados em arquivos de testes diferentes. Usei o conceito de [Build Constraints ou Build Tag](http://golang.org/pkg/go/build/#hdr-Build_Constraints) para selecionar quais testes queremos executar.

Para especificar um grupo de teste executamos o comando *go test* com o parâmetro *-tags*.

```sh
go test ./internal/pkg/commom/config -tags="unit"
```

Neste exemplo o pacote *config* possui os seguintes arquivos de teste:

- [config_test.go](internal/pkg/commom/config/config_test.go)
    ```go
    // +build test unit

    package config
    // ...
    ```

- [configEnvVarsIT_test.go](internal/pkg/commom/config/configEnvVarsIT_test.go)
    ```go
    // +build testenvvars

    package config
    // ...
    ```

- [configIT_test.go](internal/pkg/commom/config/configIT_test.go)
    ```go
    // +build test integration

    package config
    // ...
    ```

Apenas os testes do arquivo [config_test.go](internal/pkg/commom/config/config_test.go), com a build tag *"// +build test **unit**"*, serão executados pois no comando informamos *-tags="**unit**"*.

### Unitários

Testes unitários que não dependem da infra para executar, são mais rápidos, podendo conter Mock Objects conforme necessário.

Use o seguinte comando para executar os testes unitários localmente.

```sh
make test-unit
```

Estes comandos são atalhos para a execução do script [test.sh](scripts/test.sh) com parâmetro *unit*, que resulta em:

```sh
go test ./... -tags="unit" -cover -coverprofile=coverage.out
```

Para configurar um arquivo como Unit Test:

- Sufixo - *_test.go
- Build Tag - *unit*
    - Ex: // +build test unit
    - Ex: arquivo [config_test.go](internal/pkg/commom/config/config_test.go)
    - Ex: arquivo [calculator_test.go](internal/pkg/currencyexchange/calculator_test.go)

### Integração

Testes de integração dependem do [deploy da infra de testes](#Infra-Testes). Acessam os serviços de dependência sem Mock Objects. Procuramos acelerar sua execução habilitando o paralelismo com *t.Parallel()*, e para permitir isso cada teste tem seu próprio prefixo nas chaves do Redis. 

Chaves dos redis mudam de acordo com o deploy. Ex: ***[prefixo]:currency:supported***

- Em teste recebe prefixo *TestFindSupportedCurrency*.

```
TestFindSupportedCurrency:currency:supported
```

- Em produção recebe prefixo *exchange*.

```
exchange:currency:supported
```

Use o seguinte comando para executar os testes de integração localmente.

```sh
docker-compose up -d --build redis-test # Apenas uma vez caso ainda não tenha executado, mas é idempotente.
make test-integration
```

Este comando é atalho para a execução do script [test.sh](scripts/test.sh) com parâmetro *integration*, que resulta em:

```sh
go test -parallel 10 -timeout 1m30s ./... -tags="integration" -cover -coverprofile=coverage.out
```

Para configurar um arquivo como Integration Test:

- Sufixo - *IT_test.go
- Build Tag - *integration*
    - Ex: // +build test integration
    - Ex: arquivo [ratesIT_test.go](internal/pkg/currencyexchange/ratesIT_test.go)

### Unitários e Integração

Permite executar ao mesmo tempo testes de [Unidade](#Unitários) e de [Integração](#Integração). O benefício é ter maior cobertura e estatística unificada.

Os testes devem ser configurados com a build tag ***test***, sendo:

- Unitários
    ```go
    // +build test unit
    ```
- Integração 
    ```go
    // +build test integration
    ```

Use o seguinte comando para executar os testes localmente.

```sh
docker-compose up -d --build redis-test # Apenas uma vez caso ainda não tenha executado, mas é idempotente.
make test-all
```

Este comando é atalho para a execução do script [test.sh](scripts/test.sh) com parâmetro *all*, que resulta em:

```sh
go test -parallel 10 -timeout 1m30s ./... -tags="test" -cover -coverprofile=coverage.out
```

### Carga

- Execução do teste de carga containerizado

Teste de carga executa a consulta ***http://localhost:80/api/exchange?from=USD\&to=EUR\&amount=5***

```sh
make test-load-container
```

## Empacotamento

Para empacotar como imagem Docker.

- Base Alpine Linux
    ```sh
    make package # cria imagem challenge/exchange-api:latest
    ```

- Base Golang Official. Alternativa criada devido a algumas indisponibilidades percebinas no repositório apk durante desenvolvimento.
    ```sh
    make package-safe # cria imagem challenge/exchange-api:latest
    ```

## Comandos Make

Todos comandos para facilitar o desenvolvimento estão no [Makefile](Makefile).

 - Listar comandos disponíveis
```sh
make help
```
```yaml
help                           : Exibe comandos make disponíveis.
run                            : Executa aplicação empacotada em imagem Alpine Linux.
run-safe                       : Executa aplicação empacotada com imagem Golang Official(pesada).
stop                           : Pára aplicação.
build                          : Compila aplicação. Gera arquivo './exchange-api.bin'.
build-static                   : Compila aplicação com lincagem estática. Ex. 'make build-static path=./'.
wire                           : Gera/Atualiza códigos(wire_gen.go) do framework de Injeção de Dependências.
generate                       : Atualiza códigos(wire_gen.go) do framework de Injeção de Dependências.
test-unit                      : Testes de unidade
test-integration               : Testes de integração
test-all                       : Executa testes de unidade e integração.
test-unit-container            : Executa testes de unidade em ambiente containerizado.
test-integration-container     : Executa testes de integração em ambiente containerizado.
test-all-container             : Executa testes de unidade e integração em ambiente containerizado.
infra-start                    : Inicia serviços de dependência containerizados.
infra-stop                     : Interrompe serviços de dependência containerizados.
package                        : Empacota API na imagem challenge/exchange-api:latest - Alpine Linux
package-safe                   : Empacota API na imagem challenge/exchange-api:latest - Golang Official(pesada)
```

## Melhorias

### Design

- Pendente refatoração para criar *Repository*.
- Melhorar um pouco mais como o logging é feito quando existe erro. Ex: Erros de validação não devem ser logados em nível ERROR, mas em WARN, ou nem isso.

### Melhorar testes de carga

- Incluir novos testes de carga com Jmeter, que permite validar respostas da API, 
pois Apache ab é limitado quanto a isto.

### Monitoramento

- Pendentente criar configuração para monitorar e gerar alertas com Prometheus+Grafana.