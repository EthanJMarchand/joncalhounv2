package models

import "database/sql"

type Session struct {
	ID     int
	UserId int
	// Token is only set when creating a new session. When looking up a session, this will be left empty as we only store the hash of a session token and we cannot reverse it.
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	// TODO: Create the session token
	return nil, nil
}

func (ss *SessionService) LookupUser(token string) (*User, error) {
	// TODO: implement this
	return nil, nil
}
