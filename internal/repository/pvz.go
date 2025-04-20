package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/model"
)

const (
	pvzTableName = "pvz"
)

type PVZRepositoryInterface interface {
	CreatePVZ(ctx context.Context, pvzReq dto.PVZCreateRequest) (*model.PVZ, error)
}

type PVZRepository struct {
	db   *sql.DB
	psql sq.StatementBuilderType
}

func NewPVZRepository(db *sql.DB) *PVZRepository {
	return &PVZRepository{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *PVZRepository) CreatePVZ(ctx context.Context, pvzReq dto.PVZCreateRequest) (*model.PVZ, error) {
	registrationDate := time.Now()

	query, args, err := r.psql.
		Insert(pvzTableName).
		Columns("registration_date", "city").
		Values(registrationDate, pvzReq.City).
		Suffix("RETURNING id, registration_date, city").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	var createdPVZ model.PVZ
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&createdPVZ.ID, &createdPVZ.RegistrationDate, &createdPVZ.City)
	if err != nil {
		return nil, fmt.Errorf("failed to create PVZ: %w", err)
	}

	return &createdPVZ, nil
}
