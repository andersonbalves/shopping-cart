package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	log.Println("Handler started")

	return events.APIGatewayProxyResponse{Body: "Hello World!", StatusCode: http.StatusOK}, nil
}

func main() {
	lambda.Start(Handler)
}
