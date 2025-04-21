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
	GetLastProductInReception(ctx context.Context, receptionID string) (*model.Product, error)
	DeleteProduct(ctx context.Context, productID string) error
	GetProductsByReceptionID(ctx context.Context, receptionID string) ([]model.Product, error)
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

func (r *ProductRepository) GetLastProductInReception(ctx context.Context, receptionID string) (*model.Product, error) {
	query, args, err := r.psql.
		Select("id", "date_time", "type", "reception_id").
		From(productTableName).
		Where(sq.Eq{"reception_id": receptionID}).
		OrderBy("date_time DESC").
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	var product model.Product
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&product.ID, &product.DateTime, &product.Type, &product.ReceptionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no products found for this reception")
		}
		return nil, fmt.Errorf("failed to get last product: %w", err)
	}

	return &product, nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, productID string) error {
	query, args, err := r.psql.
		Delete(productTableName).
		Where(sq.Eq{"id": productID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build sql query: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func (r *ProductRepository) GetProductsByReceptionID(ctx context.Context, receptionID string) ([]model.Product, error) {
	query, args, err := r.psql.
		Select("id", "date_time", "type", "reception_id").
		From(productTableName).
		Where(sq.Eq{"reception_id": receptionID}).
		OrderBy("date_time DESC").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.ID, &product.DateTime, &product.Type, &product.ReceptionID); err != nil {
			return nil, fmt.Errorf("failed to scan product row: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating product rows: %w", err)
	}

	return products, nil
}
