#!/bin/sh

set -e

case "$1" in
    package)
        docker image build -t challenge/exchange-api -f Dockerfile.package . 
        ;;
    *)
        echo "Usage: {package}" >&2
       ;;
esac
