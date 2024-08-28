package repository

import (
	"context"
	"fmt"
	"gofi/database/entity"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestCreateSession(t *testing.T) {
	s := &entity.Session{
		UserID:    uuid.New(),
		Token:     "test token",
		ExpiredAt: time.Now(),
	}

	expectedID := uuid.New()

	tcs := []struct {
		name string
		test func(*testing.T, *SessionRepository, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, repo *SessionRepository, mock sqlmock.Sqlmock) {
				// Mock the expected query and result
				expectedCreatedAt := time.Now()
				expectedUpdatedAt := time.Now()

				mock.ExpectQuery(`INSERT INTO "session" (user_id, token, expired_at) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`).
					WithArgs(s.UserID, s.Token, s.ExpiredAt).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
						AddRow(expectedID, expectedCreatedAt, expectedUpdatedAt))

				record, err := repo.CreateSession(context.Background(), s)
				require.NoError(t, err)
				require.NotNil(t, record)
				require.Equal(t, expectedID, record.ID)
				require.Equal(t, expectedCreatedAt, record.CreatedAt)
				require.Equal(t, expectedUpdatedAt, record.UpdatedAt)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed inserting session",
			test: func(t *testing.T, repo *SessionRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO "session" (user_id, token, expired_at) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`).
					WithArgs(s.UserID, s.Token, s.ExpiredAt).
					WillReturnError(fmt.Errorf("error inserting session"))

				_, err := repo.CreateSession(context.Background(), s)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				repo := NewSessionRepository(db)
				tc.test(t, repo, mock)
			})
		})
	}
}

func TestGetSession(t *testing.T) {
	s := &entity.Session{
		UserID:    uuid.New(),
		Token:     "test token",
		ExpiredAt: time.Now(),
	}

	expectedID := uuid.New()

	tcs := []struct {
		name string
		test func(*testing.T, *SessionRepository, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, repo *SessionRepository, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "user_id", "token", "expired_at"}).
					AddRow(expectedID, s.CreatedAt, s.UpdatedAt, s.UserID, s.Token, s.ExpiredAt)

				mock.ExpectQuery(`SELECT * FROM "session" WHERE user_id=$1 AND token=$2`).
					WithArgs(expectedID, s.Token).
					WillReturnRows(rows)

				record, err := repo.GetSession(context.Background(), expectedID, s.Token)
				require.NoError(t, err)
				require.Equal(t, expectedID, record.ID)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed getting session",
			test: func(t *testing.T, repo *SessionRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT * FROM "session" WHERE user_id=$1 AND token=$2`).
					WithArgs(expectedID, s.Token).
					WillReturnError(fmt.Errorf("error getting session"))

				_, err := repo.GetSession(context.Background(), expectedID, s.Token)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				repo := NewSessionRepository(db)
				tc.test(t, repo, mock)
			})
		})
	}
}

func TestListSessions(t *testing.T) {
	s := &entity.Session{
		UserID:    uuid.New(),
		Token:     "test token",
		ExpiredAt: time.Now(),
	}

	expectedID := uuid.New()

	tcs := []struct {
		name string
		test func(*testing.T, *SessionRepository, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, repo *SessionRepository, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "user_id", "token", "expired_at"}).
					AddRow(expectedID, s.CreatedAt, s.UpdatedAt, s.UserID, s.Token, s.ExpiredAt)

				mock.ExpectQuery(`SELECT * FROM "session"`).WillReturnRows(rows)

				records, err := repo.ListSessions(context.Background())
				require.NoError(t, err)
				require.Len(t, records, 1)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed querying session",
			test: func(t *testing.T, repo *SessionRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT * FROM "session"`).WillReturnError(fmt.Errorf("error querying session"))

				_, err := repo.ListSessions(context.Background())
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				repo := NewSessionRepository(db)
				tc.test(t, repo, mock)
			})
		})
	}
}

func TestUpdateSession(t *testing.T) {
	expectedID := uuid.New()

	s := &entity.Session{
		ID:        expectedID,
		UserID:    uuid.New(),
		Token:     "test token",
		ExpiredAt: time.Now(),
	}

	ns := &entity.Session{
		ID:        expectedID,
		UserID:    uuid.New(),
		Token:     "another token",
		ExpiredAt: time.Now(),
	}

	tcs := []struct {
		name string
		test func(*testing.T, *SessionRepository, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, repo *SessionRepository, mock sqlmock.Sqlmock) {
				expectedCreatedAt := time.Now()
				expectedUpdatedAt := time.Now()

				mock.ExpectQuery(`INSERT INTO "session" (user_id, token, expired_at) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`).
					WithArgs(s.UserID, s.Token, s.ExpiredAt).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
						AddRow(expectedID, expectedCreatedAt, expectedUpdatedAt))

				cs, err := repo.CreateSession(context.Background(), s)
				require.NoError(t, err)
				require.Equal(t, expectedID, cs.ID)
				require.Equal(t, expectedCreatedAt, cs.CreatedAt)
				require.Equal(t, expectedUpdatedAt, cs.UpdatedAt)

				mock.ExpectExec(`UPDATE "session" SET user_id=?, token=?, expired_at=?, updated_at=? WHERE id=?`).
					WillReturnResult(sqlmock.NewResult(1, 1))

				up, err := repo.UpdateSession(context.Background(), ns)
				require.NoError(t, err)
				require.Equal(t, expectedID, up.ID)
				require.Equal(t, ns.UserID, up.UserID)
				require.Equal(t, ns.Token, up.Token)
				require.Equal(t, ns.ExpiredAt, up.ExpiredAt)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed updating session",
			test: func(t *testing.T, repo *SessionRepository, mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE "session" SET user_id=?, token=?, expired_at=?, updated_at=? WHERE id=?`).
					WillReturnError(fmt.Errorf("error updating session"))

				_, err := repo.UpdateSession(context.Background(), s)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				repo := NewSessionRepository(db)
				tc.test(t, repo, mock)
			})
		})
	}
}

func TestDeleteSession(t *testing.T) {
	expectedID := uuid.New()

	tcs := []struct {
		name string
		test func(*testing.T, *SessionRepository, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, repo *SessionRepository, mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM "session" WHERE id=$1`).
					WithArgs(expectedID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				err := repo.DeleteSession(context.Background(), expectedID)
				require.NoError(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed deleting session",
			test: func(t *testing.T, repo *SessionRepository, mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM "session" WHERE id=$1`).
					WithArgs(expectedID).
					WillReturnError(fmt.Errorf("error deleting session"))

				err := repo.DeleteSession(context.Background(), expectedID)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				repo := NewSessionRepository(db)
				tc.test(t, repo, mock)
			})
		})
	}
}
