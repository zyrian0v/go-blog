package views

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Login struct {
	Err error
}

var (
	username = "admin"
	password = "admin"
)

func (v Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"templates/base.html",
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
