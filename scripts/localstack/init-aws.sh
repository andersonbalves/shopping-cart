#!/bin/bash

REGION=us-east-1
STAGE=test

echo "LocalStack is running, starting setup..."

/etc/localstack/init/awsresources/lambda.sh \
    ProductsAPI \
    /lambdas/products.zip \
    products \
    ${REGION}

/etc/localstack/init/awsresources/gateway.sh \
    ProductsAPI \
    products \
    ${STAGE} \
    ${REGION} \
    "GET,POST,PUT,DELETE"

/etc/localstack/init/awsresources/dynamodb.sh \
    ShoppingCart \
    UserId \
    ProductId \
    ${REGION}

/etc/localstack/init/awsresources/lambda.sh \
    ShoppingCartAPI \
    /lambdas/shopping-cart.zip \
    shopping-cart \
    ${REGION}

/etc/localstack/init/awsresources/gateway.sh \
    ShoppingCartAPI \
    shopping-cart \
    ${STAGE} \
    ${REGION} \
    "GET,POST,PUT,DELETE"
