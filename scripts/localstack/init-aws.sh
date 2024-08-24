#!/bin/bash

API_NAME=ProductsAPI
REGION=us-east-1
STAGE=test

echo "LocalStack is running, starting setup..."

/etc/localstack/init/awsresources/dynamodb.sh \
    ShoppingCart \
    UserId \
    ProductId \
    ${REGION}

/etc/localstack/init/awsresources/lambda.sh \
    ${API_NAME} \
    /lambdas/products.zip \
    products \
    ${REGION}

/etc/localstack/init/awsresources/gateway.sh \
    ${API_NAME} \
    products \
    ${STAGE} \
    ${REGION}
