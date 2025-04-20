package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/kirillidk/pvz-service/internal/model"
)

const (
	productTableName = "products"
)

type ProductRepositoryInterface interface {
	CreateProduct(ctx context.Context, productType string, receptionID string) (*model.Product, error)
}

type ProductRepository struct {
	db   *sql.DB
	psql sq.StatementBuilderType
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, productType string, receptionID string) (*model.Product, error) {
	dateTime := time.Now()

	query, args, err := r.psql.
		Insert(productTableName).
		Columns("date_time", "type", "reception_id").
		Values(dateTime, productType, receptionID).
		Suffix("RETURNING id, date_time, type, reception_id").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	var product model.Product
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&product.ID, &product.DateTime, &product.Type, &product.ReceptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return &product, nil
}
