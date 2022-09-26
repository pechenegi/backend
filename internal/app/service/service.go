package service

import (
	"context"

	"github.com/google/uuid"
	r "github.com/pechenegi/backend/internal/app/repository"
	"github.com/pechenegi/backend/internal/pkg/models"
	"github.com/rs/zerolog"
)

type Service interface {
	SignInUser(ctx context.Context, user *models.User) (string, error)
	SignUpUser(ctx context.Context, user *models.User) (string, error)
}

type service struct {
	logger   zerolog.Logger
	userRepo r.UserRepository
}

func InitService(logger zerolog.Logger, userRepo r.UserRepository) (Service, error) {
	return &service{
		logger:   logger,
		userRepo: userRepo,
	}, nil
}

func (s *service) SignInUser(ctx context.Context, user *models.User) (string, error) {
	s.logger.Debug().Str("login", user.Login).
		Msg("checking if user with provided login exists")
	dbUser, err := s.userRepo.FindUserByLogin(ctx, user.Login)
	if err != nil {
		s.logger.Err(err).Caller().Str("login", user.Login).
			Msg("unexpected error occured while trying to find user")
		return "", err
	}

	s.logger.Debug().Str("login", user.Login).Str("id", user.ID).
		Msg("comparing provided credentials with repo credentials")
	if dbUser.Login != user.Login || dbUser.Password != user.Password {
		s.logger.Warn().Str("login", user.Login).
			Msg("incorrect credentials were provided")
		return "", ErrIncorrectCredentials
	}

	return dbUser.ID, nil
}

func (s *service) SignUpUser(ctx context.Context, user *models.User) (string, error) {
	s.logger.Debug().Str("login", user.Login).
		Msg("checking if user with provided login already exists")
	ctr, err := s.userRepo.CountUsersByLogin(ctx, user.Login)
	if err != nil {
		s.logger.Err(err).Caller().Str("login", user.Login).
			Msg("unexpected error occured while trying to count users")
		return "", err
	}
	if ctr != 0 {
		s.logger.Debug().Str("login", user.Login).Msg("user with provided login already exists")
		return "", ErrUserExists
	}

	s.logger.Debug().Str("login", user.Login).
		Msg("generating id for user")
	fillUserId(user)

	s.logger.Debug().Str("login", user.Login).Str("id", user.ID).
		Msg("creating user entry in repository")
	if err = s.userRepo.CreateUser(ctx, user); err != nil {
		s.logger.Err(err).Caller().Str("login", user.Login).Str("id", user.ID).
			Msg("unexpected error occured while trying to create new user entry")
		return "", err
	}

	return user.ID, nil
}

func fillUserId(user *models.User) {
	user.ID = uuid.NewString()
}
