# Criar dockerfile auto-documentado. 
# Self-Documenting Makefiles https://swcarpentry.github.io/make-novice/08-self-doc/index.html
.PHONY : help
help : Makefile
	@sed -n 's/^##//p' $<


## build						: Compila aplicação. Gera arquivo './exchange-api.bin'
build:
	@./scripts/compile.sh build

## build-static					: Compila aplicação com lincagem estática. Ex: make build-static path=../../
build-static:
	@./scripts/compile.sh build-static $(path)

wire:
	@./scripts/compile.sh wire
	@./scripts/compile.sh wire-testes

generate:
	@./scripts/compile.sh generate
	@./scripts/compile.sh generate-testes

test-unit:
	@./scripts/test.sh unit

test-integration:
	@./scripts/test.sh integration

## test-all					: Executa testes de unidade e integração
test-all:
	@./scripts/test.sh all

test-envvars:
	@./scripts/test.sh envvars

test-unit-container:
	@docker-compose up --build test-unit

test-integration-container:
	@docker-compose up --build test-integration

## test-all-container				: Executa testes de unidade e integração em ambiente containerizado
test-all-container:
	@docker-compose up --build test-all

test-envvars-container:
	@docker-compose up --build test-envvars

## infra-start					: Inicia serviços de dependência containerizados
infra-start:
	@docker-compose up -d --build redisdb

## infra-stop					: Interrompe serviços de dependência containerizados
infra-stop:
	@docker-compose rm -fsv redisdb

## run						: Executa aplicação
run: 
	@./scripts/deploy.sh start
	@echo "Acesse API:"
	@echo "http://`docker-compose port api 3000`/"

## stop						: Pára aplicação
stop:
	@./scripts/deploy.sh stop

package: 
	@./scripts/package.sh package
