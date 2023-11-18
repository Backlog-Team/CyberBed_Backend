package models

import "github.com/pkg/errors"

var (
	ErrNoImages          = errors.New("no images provided")
	ErrTooManyImages     = errors.New("too many images")
	ErrUserExists        = errors.New("user already exists")
	ErrNotFound          = errors.New("item not found")
	ErrNoSession         = errors.New("session not found")
	ErrIncorrectPassword = errors.New("password is incorrect")
)
