package models

import "errors"

var (
	ErrEmailTaken = errors.New("models: email address already taken")
	ErrNotFound   = errors.New("models: resource could not be found")
)
