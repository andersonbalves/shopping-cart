#!/bin/bash

TABLE_NAME=$1
PARTITION_KEY=$2
SORT_KEY=$3
REGION=$4

echo "Creating DynamoDB tables"

awslocal dynamodb create-table \
  --table-name ${TABLE_NAME} \
  --attribute-definitions AttributeName=${PARTITION_KEY},AttributeType=S AttributeName=${SORT_KEY},AttributeType=S \
  --key-schema AttributeName=${PARTITION_KEY},KeyType=HASH AttributeName=${SORT_KEY},KeyType=RANGE \
  --billing-mode PAY_PER_REQUEST \
  --endpoint-url http://localstack:4566 \
  --region ${REGION}

echo "DynamoDB table created."
