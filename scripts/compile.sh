#!/bin/sh

set -e

case "$1" in
    build)
        cd ./cmd/exchangeAPI/
        go build -ldflags \
        '-X github.com/prometheus/common/version.Version='$GIT_VERSION'
        -X github.com/prometheus/common/version.BuildDate='$DATE' 
        -X github.com/prometheus/common/version.Branch='$BRANCH' 
        -X github.com/prometheus/common/version.Revision='$GIT_REVISION'
        -X github.com/prometheus/common/version.BuildUser='$USER'' \
        -o ../../exchange-api.bin ./ 
        ;;
    build-static)
        WORKBUILD=${2:-/usr/dist}
        cd ./cmd/exchangeAPI/
        CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo \
        -ldflags \
        '-X github.com/prometheus/common/version.Version='$GIT_VERSION'
        -X github.com/prometheus/common/version.BuildDate='$DATE' 
        -X github.com/prometheus/common/version.Branch='$BRANCH' 
        -X github.com/prometheus/common/version.Revision='$GIT_REVISION'
        -X github.com/prometheus/common/version.BuildUser='$USER'' \
        -o $WORKBUILD/exchange-api.bin ./
        ;;
    wire)
        wire ./cmd/exchangeAPI/
        ;;
    wire-testes)
        wire \
        ./internal/pkg/commom/config \
        ./internal/pkg/currencyexchange \
        ./internal/pkg/infra/redis \
        ;;
    generate)
        go generate ./cmd/exchangeAPI/
        ;;
    generate-testes)
        go generate \
        ./internal/pkg/commom/config \
        ./internal/pkg/currencyexchange \
        ./internal/pkg/infra/redis \
        ;;
    *)
        echo "Usage: {build|wire}" >&2
       ;;
esac
