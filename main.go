package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/ethanjmachand/lenslocked/controllers"
	"github.com/ethanjmachand/lenslocked/views"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// main is the main function of my program.
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", controllers.StaticHandler(
		views.Must(views.Parse(filepath.Join("templates", "home.gohtml")))))

	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.Parse(filepath.Join("templates", "contact.gohtml")))))

	r.Get("/faq", controllers.StaticHandler(
		views.Must(views.Parse(filepath.Join("templates", "faq.gohtml")))))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 | Page not found", http.StatusNotFound)
	})

	fmt.Println("starting server on :3000")
	http.ListenAndServe(":3000", r)
}
