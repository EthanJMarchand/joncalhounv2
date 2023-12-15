package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/ethanjmachand/lenslocked/views"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func executeTemplate(w http.ResponseWriter, filepath string) {
	t, err := views.Parse(filepath)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w, "There was an error parsing the template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

// homeHandler is currently my own routing handler function
func homeHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "home.gohtml")
	executeTemplate(w, tplPath)
}

// contactHandler serves up the contact page when someone visits the contact route
func contactHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, filepath.Join("templates", "contact.gohtml"))
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, filepath.Join("templates", "faq.gohtml"))
}

// main is the main function of my program.
func main() {
	// Parse the template

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 | Page not found", http.StatusNotFound)
	})
	fmt.Println("starting server on :3000")
	http.ListenAndServe(":3000", r)
}
