package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	r "github.com/pechenegi/backend/internal/app/repository"
	s "github.com/pechenegi/backend/internal/app/service"
	"github.com/pechenegi/backend/internal/pkg/mocks"
	"github.com/pechenegi/backend/internal/pkg/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitHandlers(t *testing.T) {
	l := createLogger()
	r, err := r.InitUserRepository(context.Background(), l)
	require.NoError(t, err)
	s, err := s.InitService(l, r)
	require.NoError(t, err)

	h, err := InitHandlers(context.Background(), l, r, s)
	assert.NoError(t, err)
	assert.NotNil(t, h)
}

func TestPostSignIn(t *testing.T) {
	t.Run("return 200 and new user's id", func(t *testing.T) {
		h, _, ms := createHandlers(gomock.NewController(t))

		ms.EXPECT().SignInUser(
			context.Background(),
			gomock.Any(),
		).Times(1).Return("1", nil)

		req, err := createPostSignInRequest(true)
		assert.NoError(t, err)
		res := httptest.NewRecorder()
		expected := PostSignResponse{UserID: "1", Err: nil}

		h.PostSignIn(res, req)
		var actual PostSignResponse
		assert.NoError(t, json.Unmarshal(res.Body.Bytes(), &actual))

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assert.Equal(t, expected, actual)
	})

	t.Run("return 400 for nil request body", func(t *testing.T) {
		h, _, _ := createHandlers(gomock.NewController(t))

		req, err := createPostSignInRequest(false)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		h.PostSignIn(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "sign in request should have a body\n", res.Body.String())
	})

	t.Run("return 400 when cannot read body", func(t *testing.T) {
		h, _, _ := createHandlers(gomock.NewController(t))

		storedReadAll := ioReadAll
		ioReadAll = fakeReadAll
		defer restoreReadAll(storedReadAll)

		req, err := createPostSignInRequest(true)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		h.PostSignIn(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "ReadAll failed\n", res.Body.String())
	})

	t.Run("return 400 when cannot unmarshal body", func(t *testing.T) {
		h, _, _ := createHandlers(gomock.NewController(t))

		storedUnmarshal := jsonUnmarshal
		jsonUnmarshal = fakeUnmarshal
		defer restoreUnmarshal(storedUnmarshal)

		req, err := createPostSignInRequest(true)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		h.PostSignIn(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "Unmarshalling failed\n", res.Body.String())
	})

	t.Run("return 500 when service layer fails", func(t *testing.T) {
		h, _, ms := createHandlers(gomock.NewController(t))

		req, err := createPostSignInRequest(true)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		ms.EXPECT().SignInUser(
			context.Background(),
			gomock.Any(),
		).Times(1).Return("", errors.New("some err"))

		h.PostSignIn(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
		assert.Equal(t, "some err\n", res.Body.String())
	})

	t.Run("return 500 when cannot marshal response", func(t *testing.T) {
		h, _, ms := createHandlers(gomock.NewController(t))

		storedMarshal := jsonMarshal
		jsonMarshal = fakeMarshal
		defer restoreMarshal(storedMarshal)

		req, err := createPostSignInRequest(true)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		ms.EXPECT().SignInUser(
			context.Background(),
			gomock.Any(),
		).Times(1).Return("1", nil)

		h.PostSignIn(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
		assert.Equal(t, "Marshalling failed\n", res.Body.String())
	})
}

func TestPostSignUp(t *testing.T) {
	t.Run("return 201 and new user's id", func(t *testing.T) {
		h, _, ms := createHandlers(gomock.NewController(t))

		ms.EXPECT().SignUpUser(
			context.Background(),
			gomock.Any(),
		).Times(1).Return("1", nil)

		req, err := createPostSignUpRequest(true)
		assert.NoError(t, err)
		res := httptest.NewRecorder()
		expected := PostSignResponse{UserID: "1", Err: nil}

		h.PostSignUp(res, req)
		var actual PostSignResponse
		assert.NoError(t, json.Unmarshal(res.Body.Bytes(), &actual))

		assert.Equal(t, http.StatusCreated, res.Result().StatusCode)
		assert.Equal(t, expected, actual)
	})

	t.Run("return 400 for nil request body", func(t *testing.T) {
		h, _, _ := createHandlers(gomock.NewController(t))

		req, err := createPostSignUpRequest(false)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		h.PostSignUp(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "sign up request should have a body\n", res.Body.String())
	})

	t.Run("return 400 when cannot read body", func(t *testing.T) {
		h, _, _ := createHandlers(gomock.NewController(t))

		storedReadAll := ioReadAll
		ioReadAll = fakeReadAll
		defer restoreReadAll(storedReadAll)

		req, err := createPostSignUpRequest(true)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		h.PostSignUp(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "ReadAll failed\n", res.Body.String())
	})

	t.Run("return 400 when cannot unmarshal body", func(t *testing.T) {
		h, _, _ := createHandlers(gomock.NewController(t))

		storedUnmarshal := jsonUnmarshal
		jsonUnmarshal = fakeUnmarshal
		defer restoreUnmarshal(storedUnmarshal)

		req, err := createPostSignUpRequest(true)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		h.PostSignUp(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "Unmarshalling failed\n", res.Body.String())
	})

	t.Run("return 500 when service layer fails", func(t *testing.T) {
		h, _, ms := createHandlers(gomock.NewController(t))

		req, err := createPostSignUpRequest(true)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		ms.EXPECT().SignUpUser(
			context.Background(),
			gomock.Any(),
		).Times(1).Return("", errors.New("some err"))

		h.PostSignUp(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
		assert.Equal(t, "some err\n", res.Body.String())
	})

	t.Run("return 500 when cannot marshal response", func(t *testing.T) {
		h, _, ms := createHandlers(gomock.NewController(t))

		storedMarshal := jsonMarshal
		jsonMarshal = fakeMarshal
		defer restoreMarshal(storedMarshal)

		req, err := createPostSignUpRequest(true)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		ms.EXPECT().SignUpUser(
			context.Background(),
			gomock.Any(),
		).Times(1).Return("1", nil)

		h.PostSignUp(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
		assert.Equal(t, "Marshalling failed\n", res.Body.String())
	})
}

func TestGetUserDebt(t *testing.T) {
	t.Run("return current dept by user id", func(t *testing.T) {
		h, _, _ := createHandlers(gomock.NewController(t))

		req, err := createGetUserDebtRequest("1", true)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		h.GetUserDebt(res, req)
		actualDebt := new(models.DebtInfo)
		err = json.Unmarshal(res.Body.Bytes(), actualDebt)
		assert.NoError(t, err)

		expectedDebt := &models.DebtInfo{
			ID:               1,
			StartDate:        time.Date(2021, time.March, 9, 0, 0, 0, 0, time.Local),
			ContractDuration: 3,
		}
		require.Equal(t, expectedDebt, actualDebt)
	})

	t.Run("return 404 for non-existing user id", func(t *testing.T) {
		h, _, _ := createHandlers(gomock.NewController(t))

		req, err := createGetUserDebtRequest("500", true)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		h.GetUserDebt(res, req)
		assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
	})

	t.Run("return 400 when user-id header isn't set", func(t *testing.T) {
		h, _, _ := createHandlers(gomock.NewController(t))

		req, err := createGetUserDebtRequest("", false)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		h.GetUserDebt(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
	})

	t.Run("return 400 for invalid user id", func(t *testing.T) {
		h, _, _ := createHandlers(gomock.NewController(t))

		req, err := createGetUserDebtRequest("###", true)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		h.GetUserDebt(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
	})

	t.Run("return 500 when can not marshal debt", func(t *testing.T) {
		h, _, _ := createHandlers(gomock.NewController(t))

		storedMarshal := jsonMarshal
		jsonMarshal = fakeMarshal
		defer restoreMarshal(storedMarshal)

		req, err := createGetUserDebtRequest("1", true)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		h.GetUserDebt(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
	})
}

func createPostSignInRequest(hasBody bool) (*http.Request, error) {
	if hasBody {
		user := &models.User{
			Login:    "testing",
			Password: "something",
		}
		json, err := json.Marshal(user)
		if err != nil {
			return nil, err
		}
		return http.NewRequest(http.MethodPost, "/user/signin", bytes.NewBuffer(json))
	}
	return http.NewRequest(http.MethodPost, "/user/signin", nil)
}

func createPostSignUpRequest(hasBody bool) (*http.Request, error) {
	if hasBody {
		user := &models.User{
			Login:    "testing",
			Password: "something",
		}
		json, err := json.Marshal(user)
		if err != nil {
			return nil, err
		}
		return http.NewRequest(http.MethodPost, "/user/signup", bytes.NewBuffer(json))
	}
	return http.NewRequest(http.MethodPost, "/user/signup", nil)
}

func createGetUserDebtRequest(userID string, setHeader bool) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, "/user/debt", nil)
	if err != nil {
		return nil, err
	}
	if setHeader {
		req.Header.Add("user-id", userID)
	}
	return req, nil
}

func createLogger() zerolog.Logger {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02-01-2006 15:04:05 MST",
	}
	return zerolog.New(output).With().Timestamp().Logger()
}

func createHandlers(ctrl *gomock.Controller) (*handlers, *mocks.MockUserRepository, *mocks.MockService) {
	rMock := mocks.NewMockUserRepository(ctrl)
	sMock := mocks.NewMockService(ctrl)
	l := createLogger()
	return &handlers{
		logger:   l,
		userRepo: rMock,
		svc:      sMock,
	}, rMock, sMock
}

func fakeReadAll(r io.Reader) ([]byte, error)                         { return []byte{}, errors.New("ReadAll failed") }
func fakeMarshal(v interface{}) ([]byte, error)                       { return []byte{}, errors.New("Marshalling failed") }
func fakeUnmarshal(data []byte, v interface{}) error                  { return errors.New("Unmarshalling failed") }
func restoreReadAll(replace func(r io.Reader) ([]byte, error))        { ioReadAll = replace }
func restoreMarshal(replace func(v interface{}) ([]byte, error))      { jsonMarshal = replace }
func restoreUnmarshal(replace func(data []byte, v interface{}) error) { jsonUnmarshal = replace }
