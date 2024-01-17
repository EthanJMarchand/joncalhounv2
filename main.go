package main

import (
	"fmt"
	"net/http"

	"github.com/ethanjmachand/lenslocked/controllers"
	"github.com/ethanjmachand/lenslocked/models"
	"github.com/ethanjmachand/lenslocked/templates"
	"github.com/ethanjmachand/lenslocked/views"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"
)

// main is the main function of my program.
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))

	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))

	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userService := models.UserService{
		DB: db,
	}

	usersC := controllers.Users{
		UserService: &userService,
	}
	usersC.Templates.NewTPL = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.CreateUser)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Get("/users/me", usersC.CurrentUser)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 | Page not found", http.StatusNotFound)
	})

	fmt.Println("starting server on :3000")
	// 32 bit key a - Z, 0-9
	csrfKey := "g9jeH6Gc39OplfGnJKI7654FcLp521ws"
	csrfMw := csrf.Protect([]byte(csrfKey), csrf.Secure(false)) // TODO: fix this defore deploying.
	http.ListenAndServe(":3000", csrfMw(r))
}
