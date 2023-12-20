package controllers

import (
	"net/http"

	"github.com/ethanjmachand/lenslocked/views"
)

type Users struct {
	Templates struct {
		NewTPL views.Template
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	// We need a view to render
	u.Templates.NewTPL.Execute(w, nil)
}
