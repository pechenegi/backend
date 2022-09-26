package service

import "errors"

var (
	ErrIncorrectCredentials = errors.New("login and/or password are incorrect")
	ErrUserExists           = errors.New("user with provided login already exists in the system")
)
