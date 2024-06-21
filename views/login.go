package views

import (
	"errors"
	"html/template"
	"log"
	"net/http"

	"context"
	"github.com/go-session/session/v3"
)

type LogIn struct {
	Auth
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
	user, ok := store.Get("user")
	if ok {
		v.User = user.(string)
	}

	if r.Method == "POST" {
		user := r.FormValue("username")
		pass := r.FormValue("password")
		if user == username && pass == password {
			store.Set("user", user)
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
	store.Delete("user")
	err = store.Save()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", 303)
}
