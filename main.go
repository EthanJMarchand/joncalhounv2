package main

import (
	"fmt"
	"net/http"

	"github.com/ethanjmachand/lenslocked/controllers"
	"github.com/ethanjmachand/lenslocked/templates"
	"github.com/ethanjmachand/lenslocked/views"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// main is the main function of my program.
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "contact.gohtml"))))

	r.Get("/faq", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "faq.gohtml"))))

	r.Get("/login", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "login.gohtml"))))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 | Page not found", http.StatusNotFound)
	})

	fmt.Println("starting server on :3000")
	http.ListenAndServe(":3000", r)
}
