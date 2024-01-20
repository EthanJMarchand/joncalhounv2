package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/ethanjmachand/lenslocked/rand"
)

const (
	// The minimum number of bytes to be used for each session token.
	MinBytesPerToken = 32
)

type Session struct {
	ID     int
	UserId int
	// Token is only set when creating a new session. When looking up a session, this will be left empty as we only store the hash of a session token and we cannot reverse it.
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
	// BytesPerToken is used to determine how many bytes to use when generating each token session. This is value is not set, or less than the MinBytesPerToken const it be be irnored and MinBytesPerToken will be used.
	BytesPerToken int
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	session := Session{
		UserId:    userID,
		Token:     token,
		TokenHash: ss.hash(token),
	}
	row := ss.DB.QueryRow(`
		UPDATE sessions
		SET token_hash = $2
		WHERE user_id = $1
		RETURNING id;`, session.UserId, session.TokenHash)
	err = row.Scan(&session.ID)
	if err == sql.ErrNoRows {
		row := ss.DB.QueryRow(`
			INSERT INTO sessions (user_id, token_hash)
			VALUES ($1, $2)
			RETURNING id;`, session.UserId, session.TokenHash)
		err = row.Scan(&session.ID)
	}
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return &session, nil
}

func (ss *SessionService) LookupUser(token string) (*User, error) {
	tokenHash := ss.hash(token)
	var user User
	row := ss.DB.QueryRow(`
		SELECT user_id 
		FROM sessions 
		WHERE token_hash = $1;`, tokenHash)
	err := row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("lookupuser: %w", err)
	}
	row = ss.DB.QueryRow(`
		SELECT email, password_hash
		FROM users WHERE id = $1`, user.ID)
	err = row.Scan(&user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("lookupuser: %w", err)
	}
	return &user, nil
}

func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.hash(token)
	_, err := ss.DB.Exec(`
		DELETE FROM sessions
		WHERE token_Hash = $1;`, tokenHash)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

// hash is a method on a *SessionService that takes a token string, and returns a hashed token string.
func (ss *SessionService) hash(token string) string {
	tokenhash := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(tokenhash[:])
}
