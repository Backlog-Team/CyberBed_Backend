package models

import (
	"errors"
)

var (
	ErrUserExists = errors.New("user already exists")
	ErrNotFound   = errors.New("user not found")
)
