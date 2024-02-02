package models

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	DefaultResetDuration = 1 * time.Hour
)

type PasswordReset struct {
	ID        int
	UserID    int
	Token     string // Token is only set when a PasswordReset is being created.
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB            *sql.DB
	BytesPerToken int           // BytesPerToken is used to determine how many bytes to use when generating each token session. This is value is not set, or less than the MinBytesPerToken const it be be ignored and MinBytesPerToken will be used.
	Duration      time.Duration // Duration is the amount of time that a password reset is valid for.
}

func (servcice *PasswordResetService) Create(email string) (*PasswordReset, error) {
	return nil, fmt.Errorf("TODO: Implement PasswordResetService.Create")
}

func (service *PasswordResetService) Consume(token string) (*User, error) {
	return nil, fmt.Errorf("TODO: implement PasswordResetService.Consume")
}
