package repository

import (
	"encoding/json"
	"errors"
	"net/http"
	"shopping-cart/model"

	"shopping-cart/repository/dynamodb_client"
	"shopping-cart/repository/dynamodb_client/mocks"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAddItemToCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoDBClient := mocks.NewMockDynamoDBClient(ctrl)
	realGetCartForUser := GetCartForUser
	GetCartForUser = func(dbClient dynamodb_client.DynamoDBClient, userId string) (events.APIGatewayProxyResponse, error) {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: "Cart updated"}, nil
	}

	item := model.CartItem{
		UserId:      "123",
		ProductId:   "abc",
		ProductName: "Test Product",
		Quantity:    1,
	}
	itemJSON, _ := json.Marshal(item)

	mockDynamoDBClient.EXPECT().PutItem(gomock.Any(), gomock.Any()).Return(&dynamodb.PutItemOutput{}, nil).Times(1)

	request := events.APIGatewayProxyRequest{
		Body: string(itemJSON),
	}

	response, err := AddItemToCart(mockDynamoDBClient, request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "Cart updated", response.Body)
	GetCartForUser = realGetCartForUser
}

func TestAddItemToCart_UnmarshalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDynamoDBClient := mocks.NewMockDynamoDBClient(gomock.NewController(t))

	request := events.APIGatewayProxyRequest{
		Body: "{invalid-json}",
	}

	response, err := AddItemToCart(mockDynamoDBClient, request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestAddItemToCart_DynamoDBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoDBClient := mocks.NewMockDynamoDBClient(ctrl)

	item := model.CartItem{
		UserId:      "123",
		ProductId:   "abc",
		ProductName: "Test Product",
		Quantity:    1,
	}
	itemJSON, _ := json.Marshal(item)

	mockDynamoDBClient.EXPECT().PutItem(gomock.Any(), gomock.Any()).Return(nil, errors.New("DynamoDB error")).Times(1)

	request := events.APIGatewayProxyRequest{
		Body: string(itemJSON),
	}

	response, err := AddItemToCart(mockDynamoDBClient, request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	assert.Equal(t, "DynamoDB error", response.Body)
}
