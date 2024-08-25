#!/bin/bash

LAMBDA_DIRS=("products" "shopping-cart")

for LAMBDA in "${LAMBDA_DIRS[@]}"; do
  echo "Building ${LAMBDA} lambda"
  cd /app/lambdas/${LAMBDA}
  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o ${LAMBDA} main.go
  zip ${LAMBDA}.zip ${LAMBDA}
  mv ${LAMBDA}.zip ../zip/
done
