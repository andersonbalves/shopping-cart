#!/bin/bash

echo "Installing zip"
apk --no-cache add zip
echo "Building products lambda"
cd /app/lambdas/products
ls -l ./

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o products main.go

zip products.zip products
ls -l ./
mv products.zip ../zip/
