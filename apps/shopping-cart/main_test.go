package main

import (
	"net/http"

	"shopping-cart/repository"
	"shopping-cart/repository/dynamodb_client"
	"shopping-cart/repository/dynamodb_client/mocks"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	mockDbClient := &mocks.MockDynamoDBClient{}

	tests := []struct {
		name               string
		method             string
		body               string
		queryParams        map[string]string
		mockFunc           func()
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:        "GET method",
			method:      "GET",
			queryParams: map[string]string{"userId": "123"},
			mockFunc: func() {
				repository.GetCartForUser = func(client dynamodb_client.DynamoDBClient, userId string) (events.APIGatewayProxyResponse, error) {
					return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: "Cart for user 123"}, nil
				}
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       "Cart for user 123",
		},
		{
			name:   "POST method",
			method: "POST",
			body:   `{"userId":"123","productId":"abc","productName":"Product ABC","quantity":1}`,
			mockFunc: func() {
				repository.AddItemToCart = func(client dynamodb_client.DynamoDBClient, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
					return events.APIGatewayProxyResponse{StatusCode: http.StatusCreated, Body: "Item added to cart"}, nil
				}
			},
			expectedStatusCode: http.StatusCreated,
			expectedBody:       "Item added to cart",
		},
		{
			name:   "PUT method",
			method: "PUT",
			body:   `{"userId":"123","productId":"abc","productName":"Product ABC","quantity":2}`,
			mockFunc: func() {
				repository.UpdateItemInCart = func(client dynamodb_client.DynamoDBClient, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
					return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: "Item updated in cart"}, nil
				}
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       "Item updated in cart",
		},
		{
			name:        "DELETE method",
			method:      "DELETE",
			queryParams: map[string]string{"userId": "123", "productId": "abc"},
			mockFunc: func() {
				repository.RemoveItemFromCart = func(client dynamodb_client.DynamoDBClient, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
					return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: "Item removed from cart"}, nil
				}
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       "Item removed from cart",
		},
		{
			name:               "Method not allowed",
			method:             "PATCH",
			mockFunc:           func() {},
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedBody:       "Method not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			request := events.APIGatewayProxyRequest{
				HTTPMethod:            tt.method,
				Body:                  tt.body,
				QueryStringParameters: tt.queryParams,
			}

			handler := Handler(mockDbClient)
			response, err := handler(request)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, response.StatusCode)
			assert.Equal(t, tt.expectedBody, response.Body)
		})
	}
}
