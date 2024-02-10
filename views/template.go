package views

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/ethanjmachand/lenslocked/context"
	"github.com/ethanjmachand/lenslocked/models"
	"github.com/gorilla/csrf"
)

type public interface {
	Public() string
}

// Must is a helper function that takes a Template, and an error, panics if there is an error, and only returns the Template.
func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

// ParseFS takes a fs.FS, and an number of pattern strings, and returns a template, and an error.
func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl := template.New(patterns[0])
	tpl.Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", fmt.Errorf("csrfField not implemented")
		},
		"currentUser": func() (template.HTML, error) {
			return "", fmt.Errorf("currentUser not implemented")
		},
		"errors": func() []string {
			return nil
		},
	},
	)
	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}

	return Template{
		htmlTpl: tpl,
	}, nil
}

// Template type gives us the ability to write our own template functions on the *template.Template type.
type Template struct {
	htmlTpl *template.Template
}

// Execute is a method on the Template type that takes a http.ResponseWriter, a *http.Request and data, and writes to bytes.Buffer, but only to catch an error in our tpl.Funcs before setting the header for the user. At the very end, io.Copy copies the &buf and writes to w.
func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) {
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("cloning template %v", err)
		http.Error(w, "There was an error rendering the page", http.StatusInternalServerError)
		return
	}
	errMsgs := errMessages(errs...)
	tpl = tpl.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
		"currentUser": func() *models.User {
			return context.User(r.Context())
		},
		"errors": func() []string {
			return errMsgs
		},
	})
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func errMessages(errs ...error) []string {
	var msgs []string
	for _, err := range errs {
		var pubErr public
		if errors.As(err, &pubErr) {
			msgs = append(msgs, pubErr.Public())
		} else {
			fmt.Println(err)
			msgs = append(msgs, "Something went wrong.")
		}
	}
	return msgs
}
