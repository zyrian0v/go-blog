package views

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"errors"
)

type Login struct{
	Err error
}

var (
	username = "admin"
	password = "admin"
)

func (v Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

	files := []string{
		"templates/layout.html",
		"templates/login.html",
	}

	if r.Method == "POST" {
		user := r.FormValue("username")
		pass := r.FormValue("password")

		if user == username && pass == password {
			fmt.Fprint(w, "youre logged in")
			return
		} else {
			v.Err = errors.New("Wrong username or password")
		}
	}

	tmpl := template.Must(template.ParseFiles(files...))
	if err := tmpl.Execute(w, v); err != nil {
		log.Println(err)
	}
}
