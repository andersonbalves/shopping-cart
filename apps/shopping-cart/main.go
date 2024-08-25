package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"shopping-cart/repository"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type CartItem struct {
	UserId      string `json:"userId"`
	ProductId   string `json:"productId"`
	ProductName string `json:"productName"`
	Quantity    int    `json:"quantity"`
}

func createDbClient() (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}
	return dynamodb.NewFromConfig(cfg), nil
}

func Handler(dbClient *dynamodb.Client) func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		requestJson, err := json.Marshal(request)
		if err != nil {
			log.Printf("Error marshalling request: %v", err)
		} else {
			log.Printf("Received request: %s", requestJson)
		}
		defer log.Printf("Completed request: %v", request)

		switch request.HTTPMethod {
		case "GET":
			return repository.GetCartForUser(dbClient, request.QueryStringParameters["userId"])
		case "POST":
			return repository.AddItemToCart(dbClient, request)
		case "PUT":
			return repository.UpdateItemInCart(dbClient, request)
		case "DELETE":
			return repository.RemoveItemFromCart(dbClient, request)
		default:
			log.Printf("Method not allowed: %s", request.HTTPMethod)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusMethodNotAllowed,
				Body:       "Method not allowed",
			}, nil
		}
	}
}

func main() {
	dbClient, err := createDbClient()
	if err != nil {
		log.Fatalf("Failed to create DynamoDB client: %v", err)
	}
	lambda.Start(Handler(dbClient))
}
