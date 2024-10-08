#!/bin/bash

API_NAME=$1
ZIP_FILE=$2
HANDLER=$3
REGION=$4

echo "Creating Lambda function"

ENV_FILE="/lambdas/env/${HANDLER}.sh"
if [ -f "$ENV_FILE" ]; then
  source "$ENV_FILE"
  ENV_VARS="Variables={DB_HOST=${DB_HOST},DB_PORT=${DB_PORT},DB_USER=${DB_USER},DB_PASSWORD=${DB_PASSWORD},DB_NAME=${DB_NAME}}"
else
  echo "Arquivo de variáveis de ambiente não encontrado: $ENV_FILE. Continuando sem variáveis de ambiente."
  ENV_VARS="{}"
fi

echo "API_NAME: ${API_NAME}, ZIP_FILE: ${ZIP_FILE}, HANDLER: ${HANDLER}, REGION: ${REGION}, ENV_VARS: ${ENV_VARS}"

awslocal lambda create-function \
  --function-name ${API_NAME} \
  --zip-file fileb://${ZIP_FILE} \
  --handler ${HANDLER} \
  --runtime go1.x \
  --role arn:aws:iam::000000000000:role/${API_NAME} \
  --endpoint-url http://localstack:4566 \
  --environment ${ENV_VARS} \
  --region ${REGION}

if [ $? != 0 ]; then
  echo "Failed: AWS / lambda / create-function"
  exit 1
fi

echo "Lambda function created."
