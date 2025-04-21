package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kirillidk/pvz-service/internal/model"
	"github.com/kirillidk/pvz-service/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestProductRepository_CreateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	productRepo := repository.NewProductRepository(db)
	ctx := context.Background()

	tests := []struct {
		name           string
		productType    string
		receptionID    string
		mockBehavior   func()
		expectedResult *model.Product
		expectedError  error
	}{
		{
			name:        "Success",
			productType: "электроника",
			receptionID: "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "date_time", "type", "reception_id"}).
					AddRow("d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", time.Now(), "электроника", "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")

				mock.ExpectQuery(`INSERT INTO products`).
					WithArgs(sqlmock.AnyArg(), "электроника", "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11").
					WillReturnRows(rows)
			},
			expectedResult: &model.Product{
				ID:          "d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				Type:        "электроника",
				ReceptionID: "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			},
			expectedError: nil,
		},
		{
			name:        "DB Error",
			productType: "электроника",
			receptionID: "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			mockBehavior: func() {
				mock.ExpectQuery(`INSERT INTO products`).
					WithArgs(sqlmock.AnyArg(), "электроника", "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11").
					WillReturnError(errors.New("db error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("failed to create product: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			product, err := productRepo.CreateProduct(ctx, tt.productType, tt.receptionID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, product)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.ID, product.ID)
				assert.Equal(t, tt.expectedResult.Type, product.Type)
				assert.Equal(t, tt.expectedResult.ReceptionID, product.ReceptionID)
				assert.NotNil(t, product.DateTime)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestProductRepository_GetLastProductInReception(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	productRepo := repository.NewProductRepository(db)
	ctx := context.Background()
	testTime := time.Now()

	tests := []struct {
		name          string
		receptionID   string
		mockBehavior  func()
		expectedValue *model.Product
		expectedError error
	}{
		{
			name:        "Success",
			receptionID: "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "date_time", "type", "reception_id"}).
					AddRow("d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", testTime, "электроника", "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, date_time, type, reception_id FROM products WHERE reception_id = $1 ORDER BY date_time DESC LIMIT 1`)).
					WithArgs("c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11").
					WillReturnRows(rows)
			},
			expectedValue: &model.Product{
				ID:          "d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				DateTime:    testTime,
				Type:        "электроника",
				ReceptionID: "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			},
			expectedError: nil,
		},
		{
			name:        "No Products Found",
			receptionID: "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			mockBehavior: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, date_time, type, reception_id FROM products WHERE reception_id = $1 ORDER BY date_time DESC LIMIT 1`)).
					WithArgs("c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11").
					WillReturnError(sql.ErrNoRows)
			},
			expectedValue: nil,
			expectedError: errors.New("no products found for this reception"),
		},
		{
			name:        "DB Error",
			receptionID: "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			mockBehavior: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, date_time, type, reception_id FROM products WHERE reception_id = $1 ORDER BY date_time DESC LIMIT 1`)).
					WithArgs("c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11").
					WillReturnError(errors.New("db error"))
			},
			expectedValue: nil,
			expectedError: errors.New("failed to get last product: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			product, err := productRepo.GetLastProductInReception(ctx, tt.receptionID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, product)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedValue, product)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestProductRepository_DeleteProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	productRepo := repository.NewProductRepository(db)
	ctx := context.Background()

	tests := []struct {
		name          string
		productID     string
		mockBehavior  func()
		expectedError error
	}{
		{
			name:      "Success",
			productID: "d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			mockBehavior: func() {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM products WHERE id = $1`)).
					WithArgs("d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: nil,
		},
		{
			name:      "Product Not Found",
			productID: "d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			mockBehavior: func() {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM products WHERE id = $1`)).
					WithArgs("d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: errors.New("product not found"),
		},
		{
			name:      "DB Error",
			productID: "d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			mockBehavior: func() {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM products WHERE id = $1`)).
					WithArgs("d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11").
					WillReturnError(errors.New("db error"))
			},
			expectedError: errors.New("failed to delete product: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			err := productRepo.DeleteProduct(ctx, tt.productID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestProductRepository_GetProductsByReceptionID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	productRepo := repository.NewProductRepository(db)
	ctx := context.Background()
	testTime := time.Now()

	tests := []struct {
		name          string
		receptionID   string
		mockBehavior  func()
		expectedValue []model.Product
		expectedError error
	}{
		{
			name:        "Success",
			receptionID: "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "date_time", "type", "reception_id"}).
					AddRow("d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", testTime, "электроника", "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11").
					AddRow("d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12", testTime, "одежда", "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, date_time, type, reception_id FROM products WHERE reception_id = $1 ORDER BY date_time DESC`)).
					WithArgs("c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11").
					WillReturnRows(rows)
			},
			expectedValue: []model.Product{
				{
					ID:          "d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
					DateTime:    testTime,
					Type:        "электроника",
					ReceptionID: "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				},
				{
					ID:          "d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12",
					DateTime:    testTime,
					Type:        "одежда",
					ReceptionID: "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				},
			},
			expectedError: nil,
		},
		{
			name:        "Empty Result",
			receptionID: "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "date_time", "type", "reception_id"})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, date_time, type, reception_id FROM products WHERE reception_id = $1 ORDER BY date_time DESC`)).
					WithArgs("c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11").
					WillReturnRows(rows)
			},
			expectedValue: []model.Product{},
			expectedError: nil,
		},
		{
			name:        "DB Error",
			receptionID: "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			mockBehavior: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, date_time, type, reception_id FROM products WHERE reception_id = $1 ORDER BY date_time DESC`)).
					WithArgs("c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11").
					WillReturnError(errors.New("db error"))
			},
			expectedValue: nil,
			expectedError: errors.New("failed to query products: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			products, err := productRepo.GetProductsByReceptionID(ctx, tt.receptionID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, products)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedValue), len(products))
				if len(tt.expectedValue) > 0 {
					for i, expectedProduct := range tt.expectedValue {
						assert.Equal(t, expectedProduct.ID, products[i].ID)
						assert.Equal(t, expectedProduct.Type, products[i].Type)
						assert.Equal(t, expectedProduct.ReceptionID, products[i].ReceptionID)
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
