package models

import (
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// User type is the matching struct to the DB.
type User struct {
	ID           int
	Email        string
	PasswordHash string
}

// UserService type just holds the DB *sql.DB to give us access to the DB connection in main.
type UserService struct {
	DB *sql.DB
}

// Create is a method on a *UserService that takes an email, and password string and returns a *User, and an error. The function changes the email to lowercase, hashes the password, creates the user, then queries the DB to store the email and passwordhash
func (us *UserService) Create(email, password string) (*User, error) {
	email = strings.ToLower(email)
	hashedbytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	passwordHash := string(hashedbytes)

	user := User{
		Email:        email,
		PasswordHash: passwordHash,
	}
	row := us.DB.QueryRow(`
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2) RETURNING id`, email, passwordHash)
	err = row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &user, nil
}

// Authenticate is a method on a *UserService. It takes an email, and a password string, and returns a *User. First, it changes the email to lowercase, and then it queries the DB for the userId, and Password hash, and populates the user properties. Then it compares the Password hashes. If no errors, returns the User.
func (us *UserService) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(email)
	user := User{
		Email: email,
	}
	row := us.DB.QueryRow(`
		SELECT id, password_hash
		FROM users WHERE email=$1`, email)
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}
	return &user, nil
}

func (us *UserService) UpdatePassword(userID int, password string) error {
	hashedbytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("UpdatePassword: %w", err)
	}
	passwordHash := string(hashedbytes)
	_, err = us.DB.Exec(`
		UPDATE users
		SET password_hash = $2
		WHERE id = $1`, userID, passwordHash)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	return nil
}

func (us *UserService) UpdateEmail(email, newemail, password string) (*User, error) {
	newemail = strings.ToLower(newemail)
	user := User{
		Email: email,
	}
	row := us.DB.QueryRow(`
		SELECT id, password_hash
		FROM users WHERE email = $1;`, email)
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("UpdateEmail(row.scan): %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("UpdateEmail(comparehashandpassword): %w", err)
	}
	_, err = us.DB.Exec(`
	UPDATE users
	SET email = $1
	WHERE id = $2;`, newemail, user.ID)
	if err != nil {
		return nil, fmt.Errorf("updateemail: %w", err)
	}
	user.Email = newemail
	return &user, nil
}
