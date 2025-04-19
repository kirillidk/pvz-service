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
	hasOpenReception, err := r.HasOpenReception(ctx, receptionCreateReq)
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

func (r *ReceptionRepository) HasOpenReception(ctx context.Context, receptionCreateReq dto.ReceptionCreateRequest) (bool, error) {
	var exists bool

	query, _, err := r.psql.
		Select("EXISTS(SELECT 1 FROM receptions WHERE pvz_id = $1 AND status = 'in_progress')").
		ToSql()

	if err != nil {
		return false, fmt.Errorf("failed to build sql query: %w", err)
	}

	err = r.db.QueryRowContext(ctx, query, receptionCreateReq.PVZID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if open reception exists: %w", err)
	}

	return exists, nil
}
