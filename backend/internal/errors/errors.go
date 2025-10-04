package errors

import "errors"

var (
    ErrGameNotFound = errors.New("game not found")
	ErrUserNotFound = errors.New("user not found")
    ErrWrongPassword = errors.New("wrong password")
	ErrUserExists = errors.New("user exists")
)