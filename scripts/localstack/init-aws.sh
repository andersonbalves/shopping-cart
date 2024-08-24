#!/bin/bash

REGION=us-east-1
STAGE=test

echo "LocalStack is running, starting setup..."

/etc/localstack/init/awsresources/dynamodb.sh \
    ShoppingCart \
    UserId \
    ProductId \
    ${REGION}

/etc/localstack/init/awsresources/lambda.sh \
    ProductsAPI \
    /lambdas/products.zip \
    products \
    ${REGION}

/etc/localstack/init/awsresources/gateway.sh \
    ProductsAPI \
    products \
    ${STAGE} \
    ${REGION}
