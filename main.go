package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// homeHandler is currently my own routing handler function
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `
		<h1>Welcome to my awesome site!</h1>
		<a href="/contact">contact</a>
		<a href="/faq">faq</a>
	`
	fmt.Fprint(w, html)
}

// contactHandler serves up the contact page when someone visits the contact route
func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `
		<h1>Contact us.</h1>
		<a href="/">Home</a>
		<a href="/faq">faq</a>
	`
	fmt.Fprint(w, html)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `
		<h1>FAQ</h1>
		<p>  - Here are some questions and answers that some folks had for me</p>
		<ul>
			<li><strong>What is the air speed velocity of a swallow?</strong></li>
			<li>African, or american?</li>
			<li><strong>Who inventer the airplane?</strong></li>
			<li>The wright brothers</li>
			<li><strong>What is my why?</strong></li>
			<li>I do what I do becuase I feel great satisfaction and purpose from lifting people up.</li>
		</ul>
		<a href="/">Home</a>
		<a href="/contact">Contact</a>
	`
	fmt.Fprint(w, html)
}

// main is the main function of my program.
func main() {
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
