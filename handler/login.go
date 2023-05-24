package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type LoginView struct{}

var (
	username = "admin"
	password = "admin"
)

func (v LoginView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("login page")

	if r.Method == "POST" {
		user := r.FormValue("username")
		pass := r.FormValue("password")

		if user == username && pass == password {
			fmt.Fprint(w, "yourre logged in")
		} else {
			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		}
		return
	}

	files := []string{
		"views/layout.html",
		"views/login.html",
	}
	tmpl := template.Must(template.ParseFiles(files...))
	if err := tmpl.Execute(w, v); err != nil {
		log.Println(err)
	}
}
