package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pechenegi/backend/internal/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostSignUp(t *testing.T) {
	t.Run("return 201 and new user's id", func(t *testing.T) {
		req, err := createPostSignUpRequest(true)
		assert.NoError(t, err)
		res := httptest.NewRecorder()
		expected := PostSignUpResponse{UserID: "1"}

		PostSignUp(res, req)
		var actual PostSignUpResponse
		assert.NoError(t, json.Unmarshal(res.Body.Bytes(), &actual))

		assert.Equal(t, http.StatusCreated, res.Result().StatusCode)
		assert.Equal(t, expected, actual)
	})

	t.Run("return 400 for nil request body", func(t *testing.T) {
		req, err := createPostSignUpRequest(false)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		PostSignUp(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "sign up request should have a body\n", res.Body.String())
	})

	t.Run("return 400 when cannot read body", func(t *testing.T) {
		storedReadAll := ioReadAll
		ioReadAll = fakeReadAll
		defer restoreReadAll(storedReadAll)

		req, err := createPostSignUpRequest(true)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		PostSignUp(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "ReadAll failed\n", res.Body.String())
	})

	t.Run("return 400 when cannot unmarshal body", func(t *testing.T) {
		storedUnmarshal := jsonUnmarshal
		jsonUnmarshal = fakeUnmarshal
		defer restoreUnmarshal(storedUnmarshal)

		req, err := createPostSignUpRequest(true)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		PostSignUp(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "Unmarshalling failed\n", res.Body.String())
	})

	t.Run("return 500 when cannot marshal response", func(t *testing.T) {
		storedMarshal := jsonMarshal
		jsonMarshal = fakeMarshal
		defer restoreMarshal(storedMarshal)

		req, err := createPostSignUpRequest(true)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		PostSignUp(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
		assert.Equal(t, "Marshalling failed\n", res.Body.String())
	})
}

func TestGetUserDebt(t *testing.T) {
	t.Run("return current dept by user id", func(t *testing.T) {
		req, err := createGetUserDebtRequest("1", true)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		GetUserDebt(res, req)
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
		req, err := createGetUserDebtRequest("500", true)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		GetUserDebt(res, req)
		assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
	})

	t.Run("return 400 when user-id header isn't set", func(t *testing.T) {
		req, err := createGetUserDebtRequest("", false)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		GetUserDebt(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
	})

	t.Run("return 400 for invalid user id", func(t *testing.T) {
		req, err := createGetUserDebtRequest("###", true)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		GetUserDebt(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
	})

	t.Run("return 500 when can not marshal debt", func(t *testing.T) {
		storedMarshal := jsonMarshal
		jsonMarshal = fakeMarshal
		defer restoreMarshal(storedMarshal)

		req, err := createGetUserDebtRequest("1", true)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		GetUserDebt(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
	})
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

func fakeReadAll(r io.Reader) ([]byte, error)                         { return []byte{}, errors.New("ReadAll failed") }
func fakeMarshal(v interface{}) ([]byte, error)                       { return []byte{}, errors.New("Marshalling failed") }
func fakeUnmarshal(data []byte, v interface{}) error                  { return errors.New("Unmarshalling failed") }
func restoreReadAll(replace func(r io.Reader) ([]byte, error))        { ioReadAll = replace }
func restoreMarshal(replace func(v interface{}) ([]byte, error))      { jsonMarshal = replace }
func restoreUnmarshal(replace func(data []byte, v interface{}) error) { jsonUnmarshal = replace }
