package controllers

import (
	"net/http"

	"github.com/ethanjmachand/lenslocked/views"
)

func StaticHandler(tpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}

func FAQ(tpl views.Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   string
	}{
		{
			Question: "Questions 1",
			Answer:   "Answer 1",
		}, {
			Question: "Questions 2",
			Answer:   "Answer 2",
		}, {
			Question: "Question 3",
			Answer:   "Answer 3",
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, questions)
	}
}
