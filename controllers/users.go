package controllers

import (
	"fmt"
	"net/http"

	"github.com/ethanjmachand/lenslocked/models"
)

type Users struct {
	Templates struct {
		NewUser     Template
		SignIn      Template
		Currentuser Template
	}
	UserService    *models.UserService
	SessionService *models.SessionService
}

// New is the handler func for when someone visits the sign up page
func (u Users) NewUser(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.NewUser.Execute(w, r, data)
}

// CreateUser is the handler for when someone POST to /users route. It creates the session token, and creates and sets the cookie.
func (u Users) CreateUser(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Create(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		// TODO: We should show a warning about not being able to sign the user in.
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

// SignIn is the handler func for when someone visits the sign in page
func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, r, data)
}

// ProcessSignIn is the handler for when someone POST to /signin. It authenticates the password, creates the session, and then creates and sets the session cookie.
func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "users/me", http.StatusFound)
}

// CurrentUser is the handler for when you visit /users/me. It reads your session cookie, compares it to the DB, and tells you what your email is.
func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	user, err := u.SessionService.LookupUser(token)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	u.Templates.Currentuser.Execute(w, r, user)
}

// ProcessSignout is the handler for when you visit /signout. It first deletes the token, user_id, and token_id from the sessions table. Then it duplicates the session cookie with a -1 max age. It then re-directs you to the sign in page.
func (u Users) ProcessSignout(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	err = u.SessionService.Delete(token)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	deleteCookie(w, CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}
