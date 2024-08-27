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

var (
	r = &entity.Role{
		Name: "Test Role",
	}
)

func withTestDB(t *testing.T, fn func(*sqlx.DB, sqlmock.Sqlmock)) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")

	fn(db, mock)
}

func TestCreateRole(t *testing.T) {
	expectedID := uuid.New()

	tcs := []struct {
		name string
		test func(*testing.T, *RoleRepository, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, repo *RoleRepository, mock sqlmock.Sqlmock) {
				// Mock the expected query and result
				expectedCreatedAt := time.Now()
				expectedUpdatedAt := time.Now()

				mock.ExpectQuery(`INSERT INTO "role" (name) VALUES ($1) RETURNING id, created_at, updated_at`).
					WithArgs(r.Name).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
						AddRow(expectedID, expectedCreatedAt, expectedUpdatedAt))

				record, err := repo.CreateRole(context.Background(), r)
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
			name: "failed inserting role",
			test: func(t *testing.T, repo *RoleRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO "role" (name) VALUES ($1) RETURNING id, created_at, updated_at`).
					WithArgs(r.Name).
					WillReturnError(fmt.Errorf("error inserting role"))

				_, err := repo.CreateRole(context.Background(), r)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				repo := NewRoleRepository(db)
				tc.test(t, repo, mock)
			})
		})
	}
}

func TestGetRole(t *testing.T) {
	expectedID := uuid.New()

	tcs := []struct {
		name string
		test func(*testing.T, *RoleRepository, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, repo *RoleRepository, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"}).
					AddRow(expectedID, r.CreatedAt, r.UpdatedAt, r.DeletedAt, r.Name)

				mock.ExpectQuery(`SELECT * FROM "role" WHERE id=$1`).
					WithArgs(expectedID).
					WillReturnRows(rows)

				record, err := repo.GetRole(context.Background(), expectedID)
				require.NoError(t, err)
				require.Equal(t, expectedID, record.ID)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed getting role",
			test: func(t *testing.T, repo *RoleRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT * FROM "role" WHERE id=$1`).
					WithArgs(expectedID).
					WillReturnError(fmt.Errorf("error getting role"))

				_, err := repo.GetRole(context.Background(), expectedID)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				repo := NewRoleRepository(db)
				tc.test(t, repo, mock)
			})
		})
	}
}

func TestListRoles(t *testing.T) {
	expectedID := uuid.New()

	tcs := []struct {
		name string
		test func(*testing.T, *RoleRepository, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, repo *RoleRepository, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"}).
					AddRow(expectedID, r.CreatedAt, r.UpdatedAt, r.DeletedAt, r.Name)

				mock.ExpectQuery(`SELECT * FROM "role"`).WillReturnRows(rows)

				records, err := repo.ListRoles(context.Background())
				require.NoError(t, err)
				require.Len(t, records, 1)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed querying role",
			test: func(t *testing.T, repo *RoleRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT * FROM "role"`).WillReturnError(fmt.Errorf("error querying role"))

				_, err := repo.ListRoles(context.Background())
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				repo := NewRoleRepository(db)
				tc.test(t, repo, mock)
			})
		})
	}
}

func TestUpdateRole(t *testing.T) {
	expectedID := uuid.New()

	r := &entity.Role{
		ID:   expectedID,
		Name: "Test Role",
	}

	nr := &entity.Role{
		ID:   expectedID,
		Name: "Another Role",
	}

	tcs := []struct {
		name string
		test func(*testing.T, *RoleRepository, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, repo *RoleRepository, mock sqlmock.Sqlmock) {
				expectedCreatedAt := time.Now()
				expectedUpdatedAt := time.Now()

				mock.ExpectQuery(`INSERT INTO "role" (name) VALUES ($1) RETURNING id, created_at, updated_at`).
					WithArgs(r.Name).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
						AddRow(expectedID, expectedCreatedAt, expectedUpdatedAt))

				cp, err := repo.CreateRole(context.Background(), r)
				require.NoError(t, err)
				require.Equal(t, expectedID, cp.ID)
				require.Equal(t, expectedCreatedAt, cp.CreatedAt)
				require.Equal(t, expectedUpdatedAt, cp.UpdatedAt)

				mock.ExpectExec(`UPDATE "role" SET name=?, updated_at=? WHERE id=?`).
					WillReturnResult(sqlmock.NewResult(1, 1))

				up, err := repo.UpdateRole(context.Background(), nr)
				require.NoError(t, err)
				require.Equal(t, expectedID, up.ID)
				require.Equal(t, nr.Name, up.Name)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed updating role",
			test: func(t *testing.T, repo *RoleRepository, mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE "role" SET name=?, updated_at=? WHERE id=?`).
					WillReturnError(fmt.Errorf("error updating role"))

				_, err := repo.UpdateRole(context.Background(), r)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				repo := NewRoleRepository(db)
				tc.test(t, repo, mock)
			})
		})
	}
}

func TestDeleteRole(t *testing.T) {
	expectedID := uuid.New()

	tcs := []struct {
		name string
		test func(*testing.T, *RoleRepository, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, repo *RoleRepository, mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM "role" WHERE id=$1`).
					WithArgs(expectedID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				err := repo.DeleteRole(context.Background(), expectedID)
				require.NoError(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed deleting role",
			test: func(t *testing.T, repo *RoleRepository, mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM "role" WHERE id=$1`).
					WithArgs(expectedID).
					WillReturnError(fmt.Errorf("error deleting role"))

				err := repo.DeleteRole(context.Background(), expectedID)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				repo := NewRoleRepository(db)
				tc.test(t, repo, mock)
			})
		})
	}
}
