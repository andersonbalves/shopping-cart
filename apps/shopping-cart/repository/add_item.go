package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"shopping-cart/model"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func AddItemToCart(dbClient *dynamodb.Client, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Adding item to cart: %s", request.Body)

	var item model.CartItem
	err := json.Unmarshal([]byte(request.Body), &item)
	if err != nil {
		log.Printf("Error unmarshalling request body: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: err.Error()}, nil
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("ShoppingCart"),
		Item: map[string]types.AttributeValue{
			"UserId":      &types.AttributeValueMemberS{Value: item.UserId},
			"ProductId":   &types.AttributeValueMemberS{Value: item.ProductId},
			"ProductName": &types.AttributeValueMemberS{Value: item.ProductName},
			"Quantity":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", item.Quantity)},
		},
	}

	_, err = dbClient.PutItem(context.TODO(), input)
	if err != nil {
		log.Printf("Error adding item to DynamoDB: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Body: err.Error()}, nil
	}

	log.Printf("Successfully added item to cart: %v", item)
	return GetCartForUser(dbClient, item.UserId)
}
