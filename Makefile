
build:
	@./scripts/compile.sh build

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

test-all:
	@./scripts/test.sh all

test-envvars:
	@./scripts/test.sh envvars

run: 
	@./scripts/deploy.sh start
	@echo "Acesse API:"
	@echo "http://`docker-compose port api 3000`/"

stop:
	@./scripts/deploy.sh stop

package: 
	@./scripts/package.sh package

.PHONY: build 
