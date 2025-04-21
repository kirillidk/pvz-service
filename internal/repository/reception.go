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
	receptionTableName = "receptions"
)

type ReceptionRepositoryInterface interface {
	CreateReception(ctx context.Context, receptionCreateReq dto.ReceptionCreateRequest) (*model.Reception, error)
	HasOpenReception(ctx context.Context, pvzID string) (bool, error)
	GetLastOpenReception(ctx context.Context, pvzID string) (*model.Reception, error)
	CloseReception(ctx context.Context, receptionID string) (*model.Reception, error)
	GetReceptionsByPVZID(ctx context.Context, pvzID string, startDate, endDate *time.Time) ([]model.Reception, error)
}

type ReceptionRepository struct {
	db   *sql.DB
	psql sq.StatementBuilderType
}

func NewReceptionRepository(db *sql.DB) *ReceptionRepository {
	return &ReceptionRepository{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *ReceptionRepository) CreateReception(ctx context.Context, receptionCreateReq dto.ReceptionCreateRequest) (*model.Reception, error) {
	hasOpenReception, err := r.HasOpenReception(ctx, receptionCreateReq.PVZID)
	if err != nil {
		return nil, fmt.Errorf("failed to check open receptions: %w", err)
	}

	if hasOpenReception {
		return nil, fmt.Errorf("there is already an open reception for this PVZ")
	}

	dateTime := time.Now()
	query, args, err := r.psql.
		Insert(receptionTableName).
		Columns("date_time", "pvz_id", "status").
		Values(dateTime, receptionCreateReq.PVZID, "in_progress").
		Suffix("RETURNING id, date_time, pvz_id, status").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	var reception model.Reception
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&reception.ID, &reception.DateTime, &reception.PVZID, &reception.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to create reception: %w", err)
	}

	return &reception, nil
}

func (r *ReceptionRepository) HasOpenReception(ctx context.Context, pvzID string) (bool, error) {
	var exists bool

	query, _, err := r.psql.
		Select("EXISTS(SELECT 1 FROM receptions WHERE pvz_id = $1 AND status = 'in_progress')").
		ToSql()

	if err != nil {
		return false, fmt.Errorf("failed to build sql query: %w", err)
	}

	err = r.db.QueryRowContext(ctx, query, pvzID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if open reception exists: %w", err)
	}

	return exists, nil
}

func (r *ReceptionRepository) GetLastOpenReception(ctx context.Context, pvzID string) (*model.Reception, error) {
	query, args, err := r.psql.
		Select("id", "date_time", "pvz_id", "status").
		From(receptionTableName).
		Where(sq.Eq{"pvz_id": pvzID, "status": "in_progress"}).
		OrderBy("date_time DESC").
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	var reception model.Reception
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&reception.ID, &reception.DateTime, &reception.PVZID, &reception.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no open reception found for this PVZ")
		}
		return nil, fmt.Errorf("failed to get open reception: %w", err)
	}

	return &reception, nil
}

func (r *ReceptionRepository) CloseReception(ctx context.Context, receptionID string) (*model.Reception, error) {
	query, args, err := r.psql.
		Update(receptionTableName).
		Set("status", "close").
		Where(sq.Eq{"id": receptionID, "status": "in_progress"}).
		Suffix("RETURNING id, date_time, pvz_id, status").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	var reception model.Reception
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&reception.ID, &reception.DateTime, &reception.PVZID, &reception.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("reception not found or already closed")
		}
		return nil, fmt.Errorf("failed to close reception: %w", err)
	}

	return &reception, nil
}

func (r *ReceptionRepository) GetReceptionsByPVZID(ctx context.Context, pvzID string, startDate, endDate *time.Time) ([]model.Reception, error) {
	queryBuilder := r.psql.
		Select("id", "date_time", "pvz_id", "status").
		From(receptionTableName).
		Where(sq.Eq{"pvz_id": pvzID})

	if startDate != nil {
		queryBuilder = queryBuilder.Where(sq.GtOrEq{"date_time": startDate})
	}

	if endDate != nil {
		queryBuilder = queryBuilder.Where(sq.LtOrEq{"date_time": endDate})
	}

	query, args, err := queryBuilder.
		OrderBy("date_time DESC").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query receptions: %w", err)
	}
	defer rows.Close()

	var receptions []model.Reception
	for rows.Next() {
		var reception model.Reception
		if err := rows.Scan(&reception.ID, &reception.DateTime, &reception.PVZID, &reception.Status); err != nil {
			return nil, fmt.Errorf("failed to scan reception row: %w", err)
		}
		receptions = append(receptions, reception)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating reception rows: %w", err)
	}

	return receptions, nil
}
