#!/bin/bash

set -e
clear

ab -n 60000 -c 10 -m GET http://localhost:80/api/exchange?from=USD\&to=EUR\&amount=5

# Verificar status do encerramento do Teste ab, com "$?"
if [ $? -eq 0 ]
then
    exit 0
else    
    exit 1
fi
