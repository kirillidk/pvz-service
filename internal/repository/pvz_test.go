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

func TestPVZRepository_CreatePVZ(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	pvzRepo := repository.NewPVZRepository(db)
	ctx := context.Background()
	testTime := time.Now()

	tests := []struct {
		name          string
		pvzReq        dto.PVZCreateRequest
		mockBehavior  func()
		expectedPVZ   *model.PVZ
		expectedError error
	}{
		{
			name: "Success",
			pvzReq: dto.PVZCreateRequest{
				City: "Москва",
			},
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "registration_date", "city"}).
					AddRow("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", testTime, "Москва")

				mock.ExpectQuery(`INSERT INTO pvz`).
					WithArgs(sqlmock.AnyArg(), "Москва").
					WillReturnRows(rows)
			},
			expectedPVZ: &model.PVZ{
				ID:               "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				RegistrationDate: testTime,
				City:             "Москва",
			},
			expectedError: nil,
		},
		{
			name: "DB Error",
			pvzReq: dto.PVZCreateRequest{
				City: "Москва",
			},
			mockBehavior: func() {
				mock.ExpectQuery(`INSERT INTO pvz`).
					WithArgs(sqlmock.AnyArg(), "Москва").
					WillReturnError(errors.New("db error"))
			},
			expectedPVZ:   nil,
			expectedError: errors.New("failed to create PVZ: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			pvz, err := pvzRepo.CreatePVZ(ctx, tt.pvzReq)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, pvz)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPVZ.ID, pvz.ID)
				assert.Equal(t, tt.expectedPVZ.City, pvz.City)
				assert.NotNil(t, pvz.RegistrationDate)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPVZRepository_GetPVZList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	pvzRepo := repository.NewPVZRepository(db)
	ctx := context.Background()
	testTime := time.Now()

	tests := []struct {
		name          string
		filter        dto.PVZFilterQuery
		mockBehavior  func()
		expectedValue []model.PVZ
		expectedError error
	}{
		{
			name: "Success Without Date Filters",
			filter: dto.PVZFilterQuery{
				Page:  1,
				Limit: 10,
			},
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "registration_date", "city"}).
					AddRow("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", testTime, "Москва").
					AddRow("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12", testTime, "Санкт-Петербург")

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.id, p.registration_date, p.city FROM pvz p ORDER BY p.registration_date DESC LIMIT 10 OFFSET 0`)).
					WillReturnRows(rows)
			},
			expectedValue: []model.PVZ{
				{
					ID:               "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
					RegistrationDate: testTime,
					City:             "Москва",
				},
				{
					ID:               "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12",
					RegistrationDate: testTime,
					City:             "Санкт-Петербург",
				},
			},
			expectedError: nil,
		},
		{
			name: "Success With Date Filters",
			filter: dto.PVZFilterQuery{
				Page:      1,
				Limit:     10,
				StartDate: &testTime,
				EndDate:   &testTime,
			},
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "registration_date", "city"}).
					AddRow("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", testTime, "Москва")

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.id, p.registration_date, p.city FROM pvz p JOIN receptions r ON p.id = r.pvz_id WHERE r.date_time >= $1 AND r.date_time <= $2 GROUP BY p.id ORDER BY p.registration_date DESC LIMIT 10 OFFSET 0`)).
					WithArgs(testTime, testTime).
					WillReturnRows(rows)
			},
			expectedValue: []model.PVZ{
				{
					ID:               "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
					RegistrationDate: testTime,
					City:             "Москва",
				},
			},
			expectedError: nil,
		},
		{
			name: "Empty Result",
			filter: dto.PVZFilterQuery{
				Page:  1,
				Limit: 10,
			},
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "registration_date", "city"})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.id, p.registration_date, p.city FROM pvz p ORDER BY p.registration_date DESC LIMIT 10 OFFSET 0`)).
					WillReturnRows(rows)
			},
			expectedValue: []model.PVZ{},
			expectedError: nil,
		},
		{
			name: "DB Error",
			filter: dto.PVZFilterQuery{
				Page:  1,
				Limit: 10,
			},
			mockBehavior: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.id, p.registration_date, p.city FROM pvz p ORDER BY p.registration_date DESC LIMIT 10 OFFSET 0`)).
					WillReturnError(errors.New("db error"))
			},
			expectedValue: nil,
			expectedError: errors.New("failed to query pvz list: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			pvzList, err := pvzRepo.GetPVZList(ctx, tt.filter)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, pvzList)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedValue), len(pvzList))
				if len(tt.expectedValue) > 0 {
					for i, expectedPVZ := range tt.expectedValue {
						assert.Equal(t, expectedPVZ.ID, pvzList[i].ID)
						assert.Equal(t, expectedPVZ.City, pvzList[i].City)
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPVZRepository_GetPVZByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	pvzRepo := repository.NewPVZRepository(db)
	ctx := context.Background()
	testTime := time.Now()

	tests := []struct {
		name          string
		pvzID         string
		mockBehavior  func()
		expectedPVZ   *model.PVZ
		expectedError error
	}{
		{
			name:  "Success",
			pvzID: "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "registration_date", "city"}).
					AddRow("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", testTime, "Москва")

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, registration_date, city FROM pvz WHERE id = $1`)).
					WithArgs("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11").
					WillReturnRows(rows)
			},
			expectedPVZ: &model.PVZ{
				ID:               "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				RegistrationDate: testTime,
				City:             "Москва",
			},
			expectedError: nil,
		},
		{
			name:  "PVZ Not Found",
			pvzID: "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a99",
			mockBehavior: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, registration_date, city FROM pvz WHERE id = $1`)).
					WithArgs("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a99").
					WillReturnError(sql.ErrNoRows)
			},
			expectedPVZ:   nil,
			expectedError: errors.New("pvz not found"),
		},
		{
			name:  "DB Error",
			pvzID: "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			mockBehavior: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, registration_date, city FROM pvz WHERE id = $1`)).
					WithArgs("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11").
					WillReturnError(errors.New("db error"))
			},
			expectedPVZ:   nil,
			expectedError: errors.New("failed to get pvz: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			pvz, err := pvzRepo.GetPVZByID(ctx, tt.pvzID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, pvz)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPVZ.ID, pvz.ID)
				assert.Equal(t, tt.expectedPVZ.City, pvz.City)
				assert.Equal(t, tt.expectedPVZ.RegistrationDate, pvz.RegistrationDate)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
