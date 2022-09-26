package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	r "github.com/pechenegi/backend/internal/app/repository"
	s "github.com/pechenegi/backend/internal/app/service"
	m "github.com/pechenegi/backend/internal/pkg/models"
	"github.com/rs/zerolog"
)

var (
	ioReadAll     = io.ReadAll
	jsonMarshal   = json.Marshal
	jsonUnmarshal = json.Unmarshal
)

type Handlers interface {
	PostSignUp(w http.ResponseWriter, r *http.Request)
	GetUserDebt(w http.ResponseWriter, r *http.Request)
}

type handlers struct {
	logger   zerolog.Logger
	userRepo r.UserRepository
	svc      s.Service
}

func InitHandlers(ctx context.Context, logger zerolog.Logger, userRepo r.UserRepository, svc s.Service) (Handlers, error) {
	return &handlers{
		logger:   logger,
		userRepo: userRepo,
		svc:      svc,
	}, nil
}

func (h *handlers) PostSignUp(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(
			w,
			"sign up request should have a body",
			http.StatusBadRequest,
		)
		return
	}
	body, err := ioReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user := new(m.User)
	if err := jsonUnmarshal(body, user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := h.svc.SignUpUser(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json, err := jsonMarshal(PostSignUpResponse{UserID: userID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(json)
}

func (h *handlers) GetUserDebt(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromHeader(r.Context(), r.Header)
	if err != nil {
		if err.Error() == "no user-id was provided" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	debt, err := getUserDebtFromRepository(r.Context(), userID)
	if err != nil {
		if err.Error() == "non-existing user" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err.Error() == "invalid user id" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	debtJson, err := jsonMarshal(debt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(debtJson)
}

func getUserIDFromHeader(ctx context.Context, h http.Header) (string, error) {
	userID := h.Get("user-id")
	if userID == "" {
		return userID, errors.New("no user-id was provided")
	}
	return userID, nil
}

func getUserDebtFromRepository(ctx context.Context, userID string) (*m.DebtInfo, error) {
	if userID == "1" {
		return &m.DebtInfo{
			ID:               1,
			StartDate:        time.Date(2021, time.March, 9, 0, 0, 0, 0, time.Local),
			ContractDuration: 3,
		}, nil
	} else if userID == "###" {
		return nil, errors.New("invalid user id")
	} else {
		return nil, errors.New("non-existing user")
	}
}

type PostSignUpResponse struct {
	UserID string `json:"user_id,omitempty"`
}
