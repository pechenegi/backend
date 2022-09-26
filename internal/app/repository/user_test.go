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

func TestInitUserRepository(t *testing.T) {
	l := createLogger()
	repo, err := InitUserRepository(context.Background(), l)
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

	t.Run("query error", func(t *testing.T) {
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

func TestFindUserByLogin(t *testing.T) {
	login := "testing"

	t.Run("user found successfully", func(t *testing.T) {
		repo, mock, err := createRepo()
		require.NoError(t, err)

		mock.ExpectQuery("SELECT id, login, password FROM users WHERE login = (.+);").
			WithArgs(login).
			WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password"}).AddRow("1", "testing", "password"))

		user, err := repo.FindUserByLogin(context.Background(), login)
		assert.NoError(t, err)
		assert.Equal(
			t,
			&models.User{
				ID:       "1",
				Login:    "testing",
				Password: "password",
			},
			user,
		)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		repo, mock, err := createRepo()
		require.NoError(t, err)

		mock.ExpectQuery("SELECT id, login, password FROM users WHERE login = (.+);").
			WithArgs(login).
			WillReturnError(errors.New("some err"))

		user, err := repo.FindUserByLogin(context.Background(), login)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		repo, mock, err := createRepo()
		require.NoError(t, err)

		mock.ExpectQuery("SELECT id, login, password FROM users WHERE login = (.+);").
			WithArgs(login).
			WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password"}).AddRow(nil, "testing", "password"))

		user, err := repo.FindUserByLogin(context.Background(), login)
		assert.Error(t, err)
		assert.Nil(t, user)
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

func createRepo() (*userRepository, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	l := createLogger()
	return &userRepository{
		db:     db,
		logger: l,
	}, mock, nil
}
