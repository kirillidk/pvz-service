package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/model"
	"github.com/kirillidk/pvz-service/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestReceptionRepository_CreateReception(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	receptionRepo := repository.NewReceptionRepository(db)
	ctx := context.Background()
	testTime := time.Now()

	tests := []struct {
		name              string
		receptionReq      dto.ReceptionCreateRequest
		mockBehavior      func()
		expectedReception *model.Reception
		expectedError     error
	}{
		{
			name: "Success",
			receptionReq: dto.ReceptionCreateRequest{
				PVZID: "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			},
			mockBehavior: func() {
				mock.ExpectQuery(`SELECT EXISTS`).
					WithArgs("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11").
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

				rows := sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status"}).
					AddRow("c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", testTime, "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "in_progress")

				mock.ExpectQuery(`INSERT INTO receptions`).
					WithArgs(sqlmock.AnyArg(), "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "in_progress").
					WillReturnRows(rows)
			},
			expectedReception: &model.Reception{
				ID:       "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				DateTime: testTime,
				PVZID:    "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				Status:   "in_progress",
			},
			expectedError: nil,
		},
		{
			name: "Already Has Open Reception",
			receptionReq: dto.ReceptionCreateRequest{
				PVZID: "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			},
			mockBehavior: func() {
				mock.ExpectQuery(`SELECT EXISTS`).
					WithArgs("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11").
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
			},
			expectedReception: nil,
			expectedError:     errors.New("there is already an open reception for this PVZ"),
		},
		{
			name: "DB Error on Check",
			receptionReq: dto.ReceptionCreateRequest{
				PVZID: "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			},
			mockBehavior: func() {
				mock.ExpectQuery(`SELECT EXISTS`).
					WithArgs("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11").
					WillReturnError(errors.New("db error"))
			},
			expectedReception: nil,
			expectedError:     errors.New("failed to check open receptions: failed to check if open reception exists: db error"),
		},
		{
			name: "DB Error on Insert",
			receptionReq: dto.ReceptionCreateRequest{
				PVZID: "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			},
			mockBehavior: func() {
				mock.ExpectQuery(`SELECT EXISTS`).
					WithArgs("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11").
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

				mock.ExpectQuery(`INSERT INTO receptions`).
					WithArgs(sqlmock.AnyArg(), "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "in_progress").
					WillReturnError(errors.New("db error"))
			},
			expectedReception: nil,
			expectedError:     errors.New("failed to create reception: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			reception, err := receptionRepo.CreateReception(ctx, tt.receptionReq)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, reception)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedReception.ID, reception.ID)
				assert.Equal(t, tt.expectedReception.PVZID, reception.PVZID)
				assert.Equal(t, tt.expectedReception.Status, reception.Status)
				assert.NotNil(t, reception.DateTime)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestReceptionRepository_HasOpenReception(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	receptionRepo := repository.NewReceptionRepository(db)
	ctx := context.Background()
	pvzID := "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"

	tests := []struct {
		name         string
		mockBehavior func()
		expected     bool
		expectedErr  error
	}{
		{
			name: "Has Open Reception",
			mockBehavior: func() {
				mock.ExpectQuery(`SELECT EXISTS`).
					WithArgs(pvzID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
			},
			expected:    true,
			expectedErr: nil,
		},
		{
			name: "No Open Reception",
			mockBehavior: func() {
				mock.ExpectQuery(`SELECT EXISTS`).
					WithArgs(pvzID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
			},
			expected:    false,
			expectedErr: nil,
		},
		{
			name: "DB Error",
			mockBehavior: func() {
				mock.ExpectQuery(`SELECT EXISTS`).
					WithArgs(pvzID).
					WillReturnError(errors.New("db error"))
			},
			expected:    false,
			expectedErr: errors.New("failed to check if open reception exists: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			result, err := receptionRepo.HasOpenReception(ctx, pvzID)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestReceptionRepository_GetLastOpenReception(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	receptionRepo := repository.NewReceptionRepository(db)
	ctx := context.Background()
	pvzID := "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"
	testTime := time.Now()

	tests := []struct {
		name              string
		mockBehavior      func()
		expectedReception *model.Reception
		expectedError     error
	}{
		{
			name: "Success",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status"}).
					AddRow("c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", testTime, pvzID, "in_progress")

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, date_time, pvz_id, status FROM receptions WHERE`)).
					WithArgs(pvzID, "in_progress").
					WillReturnRows(rows)
			},
			expectedReception: &model.Reception{
				ID:       "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				DateTime: testTime,
				PVZID:    pvzID,
				Status:   "in_progress",
			},
			expectedError: nil,
		},
		{
			name: "No Open Reception",
			mockBehavior: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, date_time, pvz_id, status FROM receptions WHERE`)).
					WithArgs(pvzID, "in_progress").
					WillReturnError(sql.ErrNoRows)
			},
			expectedReception: nil,
			expectedError:     errors.New("no open reception found for this PVZ"),
		},
		{
			name: "DB Error",
			mockBehavior: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, date_time, pvz_id, status FROM receptions WHERE`)).
					WithArgs(pvzID, "in_progress").
					WillReturnError(errors.New("db error"))
			},
			expectedReception: nil,
			expectedError:     errors.New("failed to get open reception: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			reception, err := receptionRepo.GetLastOpenReception(ctx, pvzID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, reception)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedReception.ID, reception.ID)
				assert.Equal(t, tt.expectedReception.PVZID, reception.PVZID)
				assert.Equal(t, tt.expectedReception.Status, reception.Status)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestReceptionRepository_CloseReception(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	receptionRepo := repository.NewReceptionRepository(db)
	ctx := context.Background()
	receptionID := "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"
	testTime := time.Now()

	tests := []struct {
		name              string
		mockBehavior      func()
		expectedReception *model.Reception
		expectedError     error
	}{
		{
			name: "Success",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status"}).
					AddRow(receptionID, testTime, "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "close")

				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE receptions SET status = $1 WHERE`)).
					WithArgs("close", receptionID, "in_progress").
					WillReturnRows(rows)
			},
			expectedReception: &model.Reception{
				ID:       receptionID,
				DateTime: testTime,
				PVZID:    "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				Status:   "close",
			},
			expectedError: nil,
		},
		{
			name: "Reception Not Found or Already Closed",
			mockBehavior: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE receptions SET status = $1 WHERE`)).
					WithArgs("close", receptionID, "in_progress").
					WillReturnError(sql.ErrNoRows)
			},
			expectedReception: nil,
			expectedError:     errors.New("reception not found or already closed"),
		},
		{
			name: "DB Error",
			mockBehavior: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE receptions SET status = $1 WHERE`)).
					WithArgs("close", receptionID, "in_progress").
					WillReturnError(errors.New("db error"))
			},
			expectedReception: nil,
			expectedError:     errors.New("failed to close reception: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			reception, err := receptionRepo.CloseReception(ctx, receptionID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, reception)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedReception.ID, reception.ID)
				assert.Equal(t, tt.expectedReception.PVZID, reception.PVZID)
				assert.Equal(t, tt.expectedReception.Status, reception.Status)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestReceptionRepository_GetReceptionsByPVZID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	receptionRepo := repository.NewReceptionRepository(db)
	ctx := context.Background()
	pvzID := "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"
	testTime := time.Now()

	tests := []struct {
		name               string
		startDate          *time.Time
		endDate            *time.Time
		mockBehavior       func()
		expectedReceptions []model.Reception
		expectedError      error
	}{
		{
			name:      "Success Without Date Filters",
			startDate: nil,
			endDate:   nil,
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status"}).
					AddRow("c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", testTime, pvzID, "in_progress").
					AddRow("c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12", testTime.Add(-24*time.Hour), pvzID, "close")

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, date_time, pvz_id, status FROM receptions WHERE`)).
					WithArgs(pvzID).
					WillReturnRows(rows)
			},
			expectedReceptions: []model.Reception{
				{
					ID:       "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
					DateTime: testTime,
					PVZID:    pvzID,
					Status:   "in_progress",
				},
				{
					ID:       "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12",
					DateTime: testTime.Add(-24 * time.Hour),
					PVZID:    pvzID,
					Status:   "close",
				},
			},
			expectedError: nil,
		},
		{
			name:      "Success With Date Filters",
			startDate: &testTime,
			endDate:   &testTime,
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status"}).
					AddRow("c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", testTime, pvzID, "in_progress")

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, date_time, pvz_id, status FROM receptions WHERE`)).
					WithArgs(pvzID, testTime, testTime).
					WillReturnRows(rows)
			},
			expectedReceptions: []model.Reception{
				{
					ID:       "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
					DateTime: testTime,
					PVZID:    pvzID,
					Status:   "in_progress",
				},
			},
			expectedError: nil,
		},
		{
			name:      "Empty Result",
			startDate: nil,
			endDate:   nil,
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status"})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, date_time, pvz_id, status FROM receptions WHERE`)).
					WithArgs(pvzID).
					WillReturnRows(rows)
			},
			expectedReceptions: []model.Reception{},
			expectedError:      nil,
		},
		{
			name:      "DB Error",
			startDate: nil,
			endDate:   nil,
			mockBehavior: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, date_time, pvz_id, status FROM receptions WHERE`)).
					WithArgs(pvzID).
					WillReturnError(errors.New("db error"))
			},
			expectedReceptions: nil,
			expectedError:      errors.New("failed to query receptions: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			receptions, err := receptionRepo.GetReceptionsByPVZID(ctx, pvzID, tt.startDate, tt.endDate)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, receptions)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedReceptions), len(receptions))
				if len(tt.expectedReceptions) > 0 {
					for i, expectedReception := range tt.expectedReceptions {
						assert.Equal(t, expectedReception.ID, receptions[i].ID)
						assert.Equal(t, expectedReception.PVZID, receptions[i].PVZID)
						assert.Equal(t, expectedReception.Status, receptions[i].Status)
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
