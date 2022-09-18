package service

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	r "github.com/pechenegi/backend/internal/app/repository"
	"github.com/pechenegi/backend/internal/pkg/mocks"
	"github.com/pechenegi/backend/internal/pkg/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitService(t *testing.T) {
	l := createLogger()
	r, err := r.InitUserRepository(context.Background(), l)
	require.NoError(t, err)

	svc, err := InitService(l, r)
	assert.NoError(t, err)
	assert.NotNil(t, svc)
}

func TestSignUpUser(t *testing.T) {
	user := &models.User{
		Login:    "testing",
		Password: "testing",
	}

	t.Run("successfully sign up new user", func(t *testing.T) {
		svc, rMock := createService(gomock.NewController(t))

		count := rMock.EXPECT().CountUsersByLogin(
			context.Background(),
			gomock.Eq(user.Login),
		).Times(1).Return(0, nil)
		create := rMock.EXPECT().CreateUser(
			context.Background(),
			gomock.Not(gomock.Eq(&models.User{
				ID:       "",
				Login:    "testing",
				Password: "testing",
			})),
		).Times(1).Return(nil)
		gomock.InOrder(count, create)

		id, err := svc.SignUpUser(context.Background(), user)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})

	t.Run("user already exists", func(t *testing.T) {
		svc, rMock := createService(gomock.NewController(t))

		rMock.EXPECT().CountUsersByLogin(
			context.Background(),
			gomock.Eq(user.Login),
		).Times(1).Return(1, nil)

		id, err := svc.SignUpUser(context.Background(), user)
		assert.ErrorIs(t, err, ErrUserExists)
		assert.Empty(t, id)
	})

	t.Run("error while counting", func(t *testing.T) {
		svc, rMock := createService(gomock.NewController(t))

		rMock.EXPECT().CountUsersByLogin(
			context.Background(),
			gomock.Eq(user.Login),
		).Times(1).Return(-1, errors.New("some err"))

		id, err := svc.SignUpUser(context.Background(), user)
		assert.Error(t, err)
		assert.NotEqual(t, ErrUserExists, err)
		assert.Empty(t, id)
	})

	t.Run("error while creating", func(t *testing.T) {
		svc, rMock := createService(gomock.NewController(t))

		count := rMock.EXPECT().CountUsersByLogin(
			context.Background(),
			gomock.Eq(user.Login),
		).Times(1).Return(0, nil)
		create := rMock.EXPECT().CreateUser(
			context.Background(),
			gomock.Not(gomock.Eq(&models.User{
				ID:       "",
				Login:    "testing",
				Password: "testing",
			})),
		).Times(1).Return(errors.New("some err"))
		gomock.InOrder(count, create)

		id, err := svc.SignUpUser(context.Background(), user)
		assert.Error(t, err)
		assert.NotEqual(t, ErrUserExists, err)
		assert.Empty(t, id)
	})
}

func createLogger() zerolog.Logger {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02-01-2006 15:04:05 MST",
	}
	return zerolog.New(output).With().Timestamp().Logger()
}

func createService(ctrl *gomock.Controller) (*service, *mocks.MockUserRepository) {
	rMock := mocks.NewMockUserRepository(ctrl)
	l := createLogger()
	return &service{
		logger:   l,
		userRepo: rMock,
	}, rMock
}
