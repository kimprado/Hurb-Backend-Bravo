#!/bin/sh

set -e

case "$1" in
    start)
        docker-compose up -d --build \
        api
        ;;
    stop)
        docker-compose rm -fsv \
        api \
        redisdb
        ;;
    *)
        echo "Usage: {start|stop}" >&2
       ;;
esac
