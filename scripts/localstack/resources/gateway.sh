#!/bin/bash

API_NAME=$1
PATH_PART=$2
STAGE=$3
REGION=$4
HTTP_METHODS=$5

function fail() {
  echo $2
  exit $1
}

LAMBDA_ARN=$(awslocal lambda list-functions \
  --region ${REGION} \
  --query "Functions[?FunctionName==\`${API_NAME}\`].FunctionArn" \
  --output text)

echo "Lambda ARN: ${LAMBDA_ARN}"

echo "Creating API Gateway"

awslocal apigateway create-rest-api \
  --region ${REGION} \
  --name ${API_NAME}

[ $? == 0 ] || fail 1 "Failed: AWS / apigateway / create-rest-api"

API_ID=$(awslocal apigateway get-rest-apis \
  --region ${REGION} \
  --query "items[?name==\`${API_NAME}\`].id" \
  --output text)

PARENT_RESOURCE_ID=$(awslocal apigateway get-resources \
  --rest-api-id ${API_ID} \
  --region ${REGION} \
  --query 'items[?path==`/`].id' \
  --output text)

RESOURCE_ID=$(
  awslocal apigateway create-resource \
    --rest-api-id ${API_ID} \
    --parent-id ${PARENT_RESOURCE_ID} \
    --path-part ${PATH_PART} \
    --region ${REGION} \
    --query 'id' \
    --output text
)

[ $? == 0 ] || fail 2 "Failed: AWS / apigateway / create-resource"

IFS=',' read -ra METHODS <<<"$HTTP_METHODS"

for METHOD in "${METHODS[@]}"; do
  awslocal apigateway put-method \
    --rest-api-id ${API_ID} \
    --resource-id ${RESOURCE_ID} \
    --http-method ${METHOD} \
    --authorization-type "NONE" \
    --region ${REGION}

  [ $? == 0 ] || fail 3 "Failed: AWS / apigateway / put-method for ${METHOD}"

  awslocal apigateway put-integration \
    --rest-api-id ${API_ID} \
    --resource-id ${RESOURCE_ID} \
    --http-method ${METHOD} \
    --type AWS_PROXY \
    --integration-http-method POST \
    --uri arn:aws:apigateway:${REGION}:lambda:path/2015-03-31/functions/${LAMBDA_ARN}/invocations \
    --passthrough-behavior WHEN_NO_MATCH \
    --region ${REGION}

  [ $? == 0 ] || fail 4 "Failed: AWS / apigateway / put-integration for ${METHOD}"
done

awslocal apigateway create-deployment \
  --rest-api-id ${API_ID} \
  --stage-name ${STAGE} \
  --region ${REGION}

[ $? == 0 ] || fail 5 "Failed: AWS / apigateway / create-deployment"

ENDPOINT=http://localhost:4566/restapis/${API_ID}/${STAGE}/_user_request_/${PATH_PART}

echo "API URL: ${ENDPOINT}"
echo "API Gateway created."
