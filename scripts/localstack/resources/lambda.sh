#!/bin/bash

API_NAME=$1
ZIP_FILE=$2
HANDLER=$3
REGION=$4

echo "Creating Lambda function"

awslocal lambda create-function \
  --function-name ${API_NAME} \
  --zip-file fileb://${ZIP_FILE} \
  --handler ${HANDLER} \
  --runtime go1.x \
  --role arn:aws:iam::000000000000:role/${API_NAME} \
  --endpoint-url http://localstack:4566 \
  --region ${REGION}

if [ $? != 0 ]; then
  echo "Failed: AWS / lambda / create-function"
  exit 1
fi

echo "Lambda function created."
