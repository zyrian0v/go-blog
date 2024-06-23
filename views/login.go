package views

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"context"
	"github.com/go-session/session/v3"
)

type LogIn struct {
	Err error
}

const (
	username = "admin"
	password = "admin"
)

func (v LogIn) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if r.Method == "POST" {
		user := r.FormValue("username")
		pass := r.FormValue("password")
		if user == username && pass == password {
			store.Set("user", user)
			flash := fmt.Sprintf("Successfully logged in as %v.", user)
			store.Set("flash", flash)
			err := store.Save()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			http.Redirect(w, r, "/", 303)
			return
		} else {
			v.Err = errors.New("Wrong username or password")
		}
	}

	files := []string{
		"templates/base_no_header.html",
		"templates/login.html",
	}
	tmpl := template.Must(template.ParseFiles(files...))
	if err := tmpl.Execute(w, v); err != nil {
		log.Println(err)
	}
}

type LogOut struct{}

func (v LogOut) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	store.Delete("user")
	store.Set("flash", "Successfully logged out.")
	err = store.Save()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", 303)
}
