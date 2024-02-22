package controllers

import (
	"fmt"
	"net/http"
)

const (
	CookieSession = "session"
)

// newCookie takes the cookie name, and cookie value, and returns a *http.Cookie
func newCookie(name, value string) *http.Cookie {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}
	return &cookie
}

// setCookie takes's a http.ResponseWriter, a cookie name, and a cookie value, and sets them. This function returns nothing.
func setCookie(w http.ResponseWriter, name, value string) {
	cookie := newCookie(name, value)
	http.SetCookie(w, cookie)
}

// readCookie takes the *http.Request, and the name string, reads the cookie from the request, and returns the cookie string and an error.
func readCookie(r *http.Request, name string) (string, error) {
	c, err := r.Cookie(name)
	if err != nil {
		return "", fmt.Errorf("%s: %w", name, err)
	}
	return c.Value, nil
}

// deketeCookie takes a http.ResponseWriter, and a cookie name, writes over it with a duplicate cookie, but with the max age set to -1 telling the users browser to delete the cookie.
func deleteCookie(w http.ResponseWriter, name string) {
	cookie := newCookie(name, "")
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}
