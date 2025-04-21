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
	GetPVZList(ctx context.Context, filter dto.PVZFilterQuery) ([]model.PVZ, error)
	GetPVZByID(ctx context.Context, pvzID string) (*model.PVZ, error)
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

func (r *PVZRepository) GetPVZList(ctx context.Context, filter dto.PVZFilterQuery) ([]model.PVZ, error) {
	queryBuilder := r.psql.
		Select("p.id", "p.registration_date", "p.city").
		From(pvzTableName + " p")

	if filter.StartDate != nil || filter.EndDate != nil {
		queryBuilder = queryBuilder.Join("receptions r ON p.id = r.pvz_id")

		if filter.StartDate != nil {
			queryBuilder = queryBuilder.Where(sq.GtOrEq{"r.date_time": filter.StartDate})
		}

		if filter.EndDate != nil {
			queryBuilder = queryBuilder.Where(sq.LtOrEq{"r.date_time": filter.EndDate})
		}

		queryBuilder = queryBuilder.GroupBy("p.id")
	}

	offset := (filter.Page - 1) * filter.Limit
	queryBuilder = queryBuilder.
		Offset(uint64(offset)).
		Limit(uint64(filter.Limit)).
		OrderBy("p.registration_date DESC")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query pvz list: %w", err)
	}
	defer rows.Close()

	var pvzList []model.PVZ
	for rows.Next() {
		var pvz model.PVZ
		if err := rows.Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City); err != nil {
			return nil, fmt.Errorf("failed to scan pvz row: %w", err)
		}
		pvzList = append(pvzList, pvz)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating pvz rows: %w", err)
	}

	return pvzList, nil
}

func (r *PVZRepository) GetPVZByID(ctx context.Context, pvzID string) (*model.PVZ, error) {
	query, args, err := r.psql.
		Select("id", "registration_date", "city").
		From(pvzTableName).
		Where(sq.Eq{"id": pvzID}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	var pvz model.PVZ
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pvz not found")
		}
		return nil, fmt.Errorf("failed to get pvz: %w", err)
	}

	return &pvz, nil
}
