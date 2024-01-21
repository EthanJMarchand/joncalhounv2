package controllers

import (
	"net/http"
)

// StaticHandler takes a template, and returns a http.HandlerFunc. The func it returns runs execute, and passess no data to execute.
func StaticHandler(tpl Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, nil)
	}
}

// FAQ takes a tpl, and returns a http.HandlerFunc. It sets up the annonamous questions stuct, lays out the questions, and then executes; passing the questions []struct into execute.
func FAQ(tpl Template) http.HandlerFunc {
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
		tpl.Execute(w, r, questions)
	}
}
