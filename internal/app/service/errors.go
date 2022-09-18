package service

import "errors"

var (
	ErrUserExists = errors.New("user with provided login already exists in the system")
)
