package repository

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pechenegi/backend/internal/pkg/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitRepository(t *testing.T) {
	l := createLogger()
	repo, err := InitRepository(context.Background(), l)
	assert.NoError(t, err)
	assert.NotNil(t, repo)
}

func TestCreateUser(t *testing.T) {
	user := &models.User{
		ID:       "123",
		Login:    "testing",
		Password: "testing",
	}

	t.Run("user inserted successfully", func(t *testing.T) {
		repo, mock, err := createRepo()
		assert.NoError(t, err)

		mock.
			ExpectExec("INSERT INTO users \\(id, login, password\\) VALUES (.+);").
			WithArgs(user.ID, user.Login, user.Password).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.CreateUser(context.Background(), user)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("return error when unable to begin transaction", func(t *testing.T) {
		repo, mock, err := createRepo()
		assert.NoError(t, err)

		mock.ExpectExec("INSERT INTO users \\(id, login, password\\) VALUES (.+);").
			WithArgs(user.ID, user.Login, user.Password).
			WillReturnError(errors.New("some err"))

		err = repo.CreateUser(context.Background(), user)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCountUsersByLogin(t *testing.T) {
	login := "testing"

	t.Run("users counted successfully", func(t *testing.T) {
		repo, mock, err := createRepo()
		require.NoError(t, err)

		mock.ExpectQuery("SELECT COUNT\\(login\\) FROM users WHERE login = (.+);").
			WithArgs(login).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		ctr, err := repo.CountUsersByLogin(context.Background(), login)
		assert.NoError(t, err)
		assert.Equal(t, 1, ctr)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("users counted successfully", func(t *testing.T) {
		repo, mock, err := createRepo()
		require.NoError(t, err)

		mock.ExpectQuery("SELECT COUNT\\(login\\) FROM users WHERE login = (.+);").
			WithArgs(login).
			WillReturnError(errors.New("some err"))

		ctr, err := repo.CountUsersByLogin(context.Background(), login)
		assert.Error(t, err)
		assert.Equal(t, -1, ctr)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		repo, mock, err := createRepo()
		require.NoError(t, err)

		mock.ExpectQuery("SELECT COUNT\\(login\\) FROM users WHERE login = (.+);").
			WithArgs(login).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("not a number"))

		ctr, err := repo.CountUsersByLogin(context.Background(), login)
		assert.Error(t, err)
		assert.Equal(t, -1, ctr)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func createLogger() zerolog.Logger {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02-01-2006 15:04:05 MST",
	}
	return zerolog.New(output).With().Timestamp().Logger()
}

func createRepo() (*repository, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	l := createLogger()
	return &repository{
		db:     db,
		logger: l,
	}, mock, nil
}
