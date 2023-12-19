package controllers

import (
	"net/http"

	"github.com/ethanjmachand/lenslocked/views"
)

type Users struct {
	Templates struct {
		New views.Template
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	// We need a view to render
	u.Templates.New.Execute(w, nil)
}
