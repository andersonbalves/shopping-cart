package repository

import (
	"errors"
	"net/http"

	"shopping-cart/repository/dynamodb_client"
	"shopping-cart/repository/dynamodb_client/mocks"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRemoveItemFromCart_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoDBClient := mocks.NewMockDynamoDBClient(ctrl)
	realGetCartForUser := GetCartForUser
	GetCartForUser = func(dbClient dynamodb_client.DynamoDBClient, userId string) (events.APIGatewayProxyResponse, error) {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: "Cart updated"}, nil
	}

	mockDynamoDBClient.EXPECT().DeleteItem(gomock.Any(), gomock.Any()).Return(&dynamodb.DeleteItemOutput{}, nil).Times(1)

	request := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"userId":    "123",
			"productId": "abc",
		},
	}
	response, err := RemoveItemFromCart(mockDynamoDBClient, request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "Cart updated", response.Body)
	GetCartForUser = realGetCartForUser
}

func TestRemoveItemFromCart_DeleteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoDBClient := mocks.NewMockDynamoDBClient(ctrl)
	realGetCartForUser := GetCartForUser
	GetCartForUser = func(dbClient dynamodb_client.DynamoDBClient, userId string) (events.APIGatewayProxyResponse, error) {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: "Cart updated"}, nil
	}

	mockDynamoDBClient.EXPECT().DeleteItem(gomock.Any(), gomock.Any()).Return(nil, errors.New("DynamoDB error")).Times(1)

	request := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"userId":    "123",
			"productId": "abc",
		},
	}

	response, err := RemoveItemFromCart(mockDynamoDBClient, request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	assert.Contains(t, response.Body, "DynamoDB error")
	GetCartForUser = realGetCartForUser
}
