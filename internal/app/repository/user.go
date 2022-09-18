package repository

import (
	"context"
	"database/sql"

	"github.com/pechenegi/backend/internal/pkg/models"
	"github.com/rs/zerolog"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	CountUsersByLogin(ctx context.Context, login string) (int, error)
}

type userRepository struct {
	db     *sql.DB
	logger zerolog.Logger
}

func InitUserRepository(ctx context.Context, logger zerolog.Logger) (UserRepository, error) {
	return &userRepository{
		db:     &sql.DB{},
		logger: logger,
	}, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	r.logger.Debug().Str("id", user.ID).Str("login", user.Login).
		Msg("inserting new user in the database")

	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO users (id, login, password) VALUES ($1, $2, $3);",
		user.ID,
		user.Login,
		user.Password,
	)
	if err != nil {
		r.logger.Err(err).Caller().Msg("unable to insert new user")
		return err
	}

	r.logger.Debug().Str("id", user.ID).Str("login", user.Login).
		Msg("new user was successfully inserted in the database")
	return nil
}

func (r *userRepository) CountUsersByLogin(ctx context.Context, login string) (int, error) {
	r.logger.Debug().Str("login", login).
		Msg("counting users with provided login")

	row := r.db.QueryRowContext(
		ctx,
		"SELECT COUNT(login) FROM users WHERE login = $1;",
		login,
	)
	if row.Err() != nil {
		r.logger.Err(row.Err()).Caller().Str("login", login).
			Msg("unable to execute count query")
		return -1, row.Err()
	}
	var ctr int
	if err := row.Scan(&ctr); err != nil {
		r.logger.Err(err).Caller().Str("login", login).
			Msg("unable to scan count query result")
		return -1, err
	}

	r.logger.Debug().Str("login", login).Int("counter", ctr).
		Msg("users with provided login were successfully counted")
	return ctr, nil
}
