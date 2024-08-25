package repository

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func RemoveItemFromCart(dbClient *dynamodb.Client, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userId := request.QueryStringParameters["userId"]
	productId := request.QueryStringParameters["productId"]
	log.Printf("Removing item from cart: userId=%s, productId=%s", userId, productId)

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("ShoppingCart"),
		Key: map[string]types.AttributeValue{
			"UserId":    &types.AttributeValueMemberS{Value: userId},
			"ProductId": &types.AttributeValueMemberS{Value: productId},
		},
	}

	_, err := dbClient.DeleteItem(context.TODO(), input)
	if err != nil {
		log.Printf("Error removing item from DynamoDB: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Body: err.Error()}, nil
	}

	log.Printf("Successfully removed item from cart: userId=%s, productId=%s", userId, productId)
	return GetCartForUser(dbClient, userId)
}
