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

func TestUpdateItemInCart_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoDBClient := mocks.NewMockDynamoDBClient(ctrl)
	realGetCartForUser := GetCartForUser
	GetCartForUser = func(dbClient dynamodb_client.DynamoDBClient, userId string) (events.APIGatewayProxyResponse, error) {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: "Cart updated"}, nil
	}

	item := model.CartItem{
		UserId:    "123",
		ProductId: "abc",
		Quantity:  2,
	}
	body, _ := json.Marshal(item)

	mockDynamoDBClient.EXPECT().UpdateItem(gomock.Any(), gomock.Any()).Return(&dynamodb.UpdateItemOutput{}, nil).Times(1)

	request := events.APIGatewayProxyRequest{Body: string(body)}
	response, err := UpdateItemInCart(mockDynamoDBClient, request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "Cart updated", response.Body)
	GetCartForUser = realGetCartForUser
}

func TestUpdateItemInCart_UnmarshalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoDBClient := mocks.NewMockDynamoDBClient(ctrl)
	realGetCartForUser := GetCartForUser
	GetCartForUser = func(dbClient dynamodb_client.DynamoDBClient, userId string) (events.APIGatewayProxyResponse, error) {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: "Cart updated"}, nil
	}

	request := events.APIGatewayProxyRequest{Body: "invalid body"}
	response, err := UpdateItemInCart(mockDynamoDBClient, request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Contains(t, response.Body, "invalid character 'i' looking for beginning of value")
	GetCartForUser = realGetCartForUser
}

func TestUpdateItemInCart_UpdateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoDBClient := mocks.NewMockDynamoDBClient(ctrl)
	realGetCartForUser := GetCartForUser
	GetCartForUser = func(dbClient dynamodb_client.DynamoDBClient, userId string) (events.APIGatewayProxyResponse, error) {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: "Cart updated"}, nil
	}

	item := model.CartItem{
		UserId:    "123",
		ProductId: "abc",
		Quantity:  2,
	}
	body, _ := json.Marshal(item)

	mockDynamoDBClient.EXPECT().UpdateItem(gomock.Any(), gomock.Any()).Return(nil, errors.New("DynamoDB error")).Times(1)

	request := events.APIGatewayProxyRequest{Body: string(body)}
	response, err := UpdateItemInCart(mockDynamoDBClient, request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	assert.Contains(t, response.Body, "DynamoDB error")
	GetCartForUser = realGetCartForUser
}
