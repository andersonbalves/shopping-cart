package repository

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"shopping-cart/model"
	"shopping-cart/repository/dynamodb_client"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var GetCartForUser = func(dbClient dynamodb_client.DynamoDBClient, userId string) (events.APIGatewayProxyResponse, error) {
	log.Printf("Getting cart for user: %s", userId)

	input := &dynamodb.QueryInput{
		TableName:              aws.String("ShoppingCart"),
		KeyConditionExpression: aws.String("UserId = :userId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId": &types.AttributeValueMemberS{Value: userId},
		},
	}

	result, err := dbClient.Query(context.TODO(), input)
	if err != nil {
		log.Printf("Error querying DynamoDB: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Body: err.Error()}, nil
	}

	items := []model.CartItem{}
	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		log.Printf("Error unmarshalling result: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Body: err.Error()}, nil
	}

	body, err := json.Marshal(items)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Body: err.Error()}, nil
	}

	log.Printf("Successfully retrieved cart for user: %s", userId)
	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: string(body)}, nil
}
