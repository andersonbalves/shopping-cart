package main

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestFetchProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(1, "Product1", 100.0).
		AddRow(2, "Product2", 200.0)
	mock.ExpectQuery("SELECT id, name, price FROM products").WillReturnRows(rows)

	products, err := fetchProducts(db)

	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, products[0].ID, 1)
	assert.Equal(t, products[0].Name, "Product1")
	assert.Equal(t, products[0].Price, 100.0)
	assert.Equal(t, products[1].ID, 2)
	assert.Equal(t, products[1].Name, "Product2")
	assert.Equal(t, products[1].Price, 200.0)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetDBConnection(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_USER", "user")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("DB_NAME", "testdb")

	db, err := getDBConnection()

	assert.NoError(t, err)
	assert.NotNil(t, db)

	db.Close()
}

func TestHandler(t *testing.T) {
	ctx := context.Background()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(1, "Product1", 100.0).
		AddRow(2, "Product2", 200.0)
	mock.ExpectQuery("SELECT id, name, price FROM products").WillReturnRows(rows)

	mockGetDBConnection := func() (*sql.DB, error) {
		return db, nil
	}

	response, err := Handler(ctx, mockGetDBConnection)

	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.Contains(t, response.Body, "Product1")
	assert.Contains(t, response.Body, "Product2")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
