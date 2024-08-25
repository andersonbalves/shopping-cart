package repository

import (
	"errors"
	"net/http"
	"shopping-cart/repository/dynamodb_client/mocks"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetCartForUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoDBClient := mocks.NewMockDynamoDBClient(ctrl)

	mockDynamoDBClient.EXPECT().Query(gomock.Any(), gomock.Any()).Return(&dynamodb.QueryOutput{
		Items: []map[string]types.AttributeValue{
			{
				"UserId":      &types.AttributeValueMemberS{Value: "123"},
				"ProductId":   &types.AttributeValueMemberS{Value: "abc"},
				"ProductName": &types.AttributeValueMemberS{Value: "Test Product"},
				"Quantity":    &types.AttributeValueMemberN{Value: "1"},
			},
		},
	}, nil).Times(1)

	response, err := GetCartForUser(mockDynamoDBClient, "123")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	expectedBody := `[{"userId":"123","productId":"abc","productName":"Test Product","quantity":1}]`
	assert.JSONEq(t, expectedBody, response.Body)
}

func TestGetCartForUser_QueryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoDBClient := mocks.NewMockDynamoDBClient(ctrl)

	mockDynamoDBClient.EXPECT().Query(gomock.Any(), gomock.Any()).Return(nil, errors.New("DynamoDB error")).Times(1)

	response, err := GetCartForUser(mockDynamoDBClient, "123")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	assert.Equal(t, "DynamoDB error", response.Body)
}
