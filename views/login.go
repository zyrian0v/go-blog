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

var (
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
			store.Set("auth", true)
			err := store.Save()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			fmt.Fprint(w, "you are logged in")
			return
		} else {
			v.Err = errors.New("Wrong username or password")
		}
	}

	files := []string{
		"templates/base.html",
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
	store.Delete("auth")
	err = store.Save()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprint(w, "you are logged out")
}
