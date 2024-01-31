package main

import (
	"fmt"
	"net/http"

	"github.com/ethanjmachand/lenslocked/controllers"
	"github.com/ethanjmachand/lenslocked/migrations"
	"github.com/ethanjmachand/lenslocked/models"
	"github.com/ethanjmachand/lenslocked/templates"
	"github.com/ethanjmachand/lenslocked/views"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"
)

func main() {
	// Setup DB.
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Setup Services
	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}

	// Setup middleware
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	csrfKey := "g9jeH6Gc39OplfGnJKI7654FcLp521ws"
	csrfMw := csrf.Protect([]byte(csrfKey), csrf.Secure(false)) // TODO: fix this defore deploying.

	// Setup controllers
	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	usersC.Templates.NewUser = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	usersC.Templates.Currentuser = views.Must(views.ParseFS(templates.FS, "currentuser.gohtml", "tailwind.gohtml"))

	// Setup our router and routes
	r := chi.NewRouter()
	r.Use(csrfMw)
	r.Use(umw.SetUser)
	r.Use(middleware.Logger)
	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))
	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))
	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))
	r.Get("/signup", usersC.NewUser)
	r.Post("/users", usersC.CreateUser)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignout)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})
	// r.Get("/users/me", usersC.CurrentUser)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 | Page not found", http.StatusNotFound)
	})

	fmt.Println("starting server on :3000")
	http.ListenAndServe(":3000", r)
}
