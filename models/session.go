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

// Session type is all the properties needed to setup a session for a user.
type Session struct {
	ID        int
	UserId    int
	Token     string // Token is only set when creating a new session. When looking up a session, this will be left empty as we only store the hash of a session token and we cannot reverse it.
	TokenHash string
}

// SessionService type holds out *sql.DB so that we have access to the DB connection when using it in main.
type SessionService struct {
	DB            *sql.DB
	BytesPerToken int // BytesPerToken is used to determine how many bytes to use when generating each token session. This is value is not set, or less than the MinBytesPerToken const it be be ignored and MinBytesPerToken will be used.
}

// Create is a method on a *SessionService that takes a userID int, and returns a *Session and an error. This will use the MinBytesPerToken const if the BytesPertoken properties is set to less than 32, or is not set at all. Create pushes the Session to the DB.
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
		INSERT INTO sessions (user_id, token_hash)
		VALUES($1, $2) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2
		RETURNING id;`, session.UserId, session.TokenHash)
	err = row.Scan(&session.ID)
	// if err == sql.ErrNoRows {
	// 	row := ss.DB.QueryRow(`
	// 		INSERT INTO sessions (user_id, token_hash)
	// 		VALUES ($1, $2)
	// 		RETURNING id;`, session.UserId, session.TokenHash)
	// 	err = row.Scan(&session.ID)
	// }
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return &session, nil
}

// LookupUser is a method on the *SessionService type that takes the token string, hashes it, queries the DB with the hashedtoken string, and returns the user, if exists.
func (ss *SessionService) LookupUser(token string) (*User, error) {
	tokenHash := ss.hash(token)
	var user User
	row := ss.DB.QueryRow(`
		SELECT users.id, users.email, users.password_hash
		FROM sessions
		JOIN users ON users.id = sessions.user_id
		WHERE sessions.token_hash = $1;`, tokenHash)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("lookupuser: %w", err)
	}
	return &user, nil
}

// Delete is a method on the *SessionService type that takes the token string, hashes it, and then deletes the data entry from the DB.
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
