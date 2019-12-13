#!/bin/sh

set -e

case "$1" in
    start)
        docker-compose up -d --build \
        api
        ;;
    start-safe)
        docker-compose up -d --build \
        api-safe
        ;;
    stop)
        docker-compose rm -fsv \
        api \
        api-safe \
        redisdb
        ;;
    *)
        echo "Usage: {start|stop}" >&2
       ;;
esac
