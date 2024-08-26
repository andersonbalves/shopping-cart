package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"shopping-cart/model"
	"shopping-cart/repository/dynamodb_client"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var UpdateItemInCart = func(dbClient dynamodb_client.DynamoDBClient, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Updating item in cart: %s", request.Body)

	var item model.CartItem
	err := json.Unmarshal([]byte(request.Body), &item)
	if err != nil {
		log.Printf("Error unmarshalling request body: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: err.Error()}, nil
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("ShoppingCart"),
		Key: map[string]types.AttributeValue{
			"UserId":    &types.AttributeValueMemberS{Value: item.UserId},
			"ProductId": &types.AttributeValueMemberS{Value: item.ProductId},
		},
		UpdateExpression: aws.String("set Quantity = :quantity"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":quantity": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", item.Quantity)},
		},
	}

	_, err = dbClient.UpdateItem(context.TODO(), input)
	if err != nil {
		log.Printf("Error updating item in DynamoDB: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Body: err.Error()}, nil
	}

	log.Printf("Successfully updated item in cart: %v", item)
	return GetCartForUser(dbClient, item.UserId)
}
