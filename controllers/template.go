package controllers

import "net/http"

// Template interface makes it so that the controllers package does not have to import the views package. Instead anything with the Execute(w, r, data) signature can be used, even external packages.
type Template interface {
	Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error)
}
