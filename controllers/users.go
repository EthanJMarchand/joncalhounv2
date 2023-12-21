package controllers

import (
	"fmt"
	"net/http"
)

type Users struct {
	Templates struct {
		NewTPL Template
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	// We need a view to render
	u.Templates.NewTPL.Execute(w, nil)
}

func (u Users) CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Email:", r.FormValue("email"), "\n")
	fmt.Fprint(w, "Password:", r.FormValue("password"))
}
