package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ethanjmachand/lenslocked/context"
	"github.com/ethanjmachand/lenslocked/errors"
	"github.com/ethanjmachand/lenslocked/models"
)

type Users struct {
	Templates struct {
		NewUser        Template
		SignIn         Template
		Currentuser    Template
		ForgotPassword Template
		CheckYourEmail Template
		ResetPassword  Template
		Account        Template
		UpdateEmail    Template
	}
	UserService           *models.UserService
	SessionService        *models.SessionService
	PasswordResestService *models.PasswordResetService
	EmailService          *models.EmailService
}

// ---------------------------------------------------------------------------
// New is the handler func for when someone visits the sign up page
func (u Users) NewUser(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.NewUser.Execute(w, r, data)
}

// ---------------------------------------------------------------------------
// CreateUser is the handler for when someone POST to /users route. It creates the session token, and creates and sets the cookie.
func (u Users) CreateUser(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	user, err := u.UserService.Create(data.Email, data.Password)
	if err != nil {
		if errors.Is(err, models.ErrEmailTaken) {
			err = errors.Public(err, "That email is already associated with an account.")
		}
		u.Templates.NewUser.Execute(w, r, data, err)
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
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

// ---------------------------------------------------------------------------
// SignIn is the handler func for when someone visits the sign in page
func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, r, data)
}

// ---------------------------------------------------------------------------
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
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

// ---------------------------------------------------------------------------
// CurrentUser is the handler for when you visit /users/me. It reads your session cookie, compares it to the DB, and tells you what your email is.
func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	u.Templates.Currentuser.Execute(w, r, user)
}

// ---------------------------------------------------------------------------
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

// ---------------------------------------------------------------------------
// ForgotPassword is a method on the User type that is also a handler func. This Handler func would handle a POST request, and send the person trying to login an email to help them create a new password.
func (u Users) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.ForgotPassword.Execute(w, r, data)
}

// ---------------------------------------------------------------------------
// ProcessForgotPassword is a reciever function on the User, and it is also a handlerfunc.
func (u Users) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	pwReset, err := u.PasswordResestService.Create(data.Email)
	if err != nil {
		// TODO: Handle other cases in the future for example if a user does not exist with that email.
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	vals := url.Values{
		"token": {pwReset.Token},
	}
	resetURL := "https://localhost:3000/reset-pw?" + vals.Encode()
	err = u.EmailService.ForgotPassword(data.Email, resetURL)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	u.Templates.CheckYourEmail.Execute(w, r, data)
}

func (u Users) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token string
	}
	data.Token = r.FormValue("token")
	u.Templates.ResetPassword.Execute(w, r, data)
}

func (u Users) ProcessResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token    string
		Password string
	}
	data.Token = r.FormValue("token")
	data.Password = r.FormValue("password")

	user, err := u.PasswordResestService.Consume(data.Token)
	if err != nil {
		fmt.Println(err)
		// TODO: Distinguish between types of errors.
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = u.UserService.UpdatePassword(user.ID, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Sign the user in now that their password has been reset.
	// Any errors from this point onward should redirect the user to the sign in page.
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) Account(w http.ResponseWriter, r *http.Request) {
	u.Templates.Account.Execute(w, r, nil)
}

func (u Users) UpdateEmail(w http.ResponseWriter, r *http.Request) {
	u.Templates.UpdateEmail.Execute(w, r, nil)
}

func (u Users) ProcessUpdateEmail(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Newemail string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Newemail = r.FormValue("newemail")
	data.Password = r.FormValue("password")
	user, err := u.UserService.UpdateEmail(data.Email, data.Newemail, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	fmt.Println(user)
	http.Redirect(w, r, "/users/me/account", http.StatusFound)
}

// ---------------------------------------------------------------------------
type UserMiddleware struct {
	SessionService *models.SessionService
}

// ---------------------------------------------------------------------------
func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := readCookie(r, CookieSession)
		if err != nil {
			fmt.Println("SetUser: readCookie error")
			next.ServeHTTP(w, r)
			return
		}
		user, err := umw.SessionService.LookupUser(token)
		if err != nil {
			fmt.Println("SetUser: LookupUser error")
			next.ServeHTTP(w, r)
			return
		}
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// ---------------------------------------------------------------------------
func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
