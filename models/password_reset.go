package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/ethanjmachand/lenslocked/rand"
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

func (service *PasswordResetService) Create(email string) (*PasswordReset, error) {
	// verify we have a valid email address, and get their ID
	email = strings.ToLower(email)
	var userID int
	row := service.DB.QueryRow(`
	SELECT id
	FROM Users
	WHERE email = $1;`, email)
	err := row.Scan(&userID)
	if err != nil {
		// TODO: Consider returning a specific err when the user does not exist
		return nil, fmt.Errorf("create: %w", err)
	}
	// Build the PasswordReset
	bytesPerToken := service.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

}

func (service *PasswordResetService) Consume(token string) (*User, error) {
	return nil, fmt.Errorf("TODO: implement PasswordResetService.Consume")
}

func (service *PasswordResetService) hash(token string) string {
	tokenhash := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(tokenhash[:])
}
