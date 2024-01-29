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

// main is the main function of my program.
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))

	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))

	// Here, we setup out PostgresConfig properties, and open the Database connection. We have to remember to defer db.Close().
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

	// Here, we pass the DB connection to our userService, and our sessionService.
	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}

	// here, we setup out controller to have access to the UserService, and SessionService.
	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	usersC.Templates.NewUser = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	usersC.Templates.Currentuser = views.Must(views.ParseFS(templates.FS, "currentuser.gohtml", "tailwind.gohtml"))
	r.Get("/signup", usersC.NewUser)
	r.Post("/users", usersC.CreateUser)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignout)
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
