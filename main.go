package main

import (
	"log"
	"net/http"

	"blog/db"
	"blog/views"
)

func main() {
	// Database
	db.InitializeHandle()
	db.ApplySchema()

	// Routes
	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	root := views.Index{}
	http.Handle("/", root)

	http.Handle("/login/", views.Login{})
	http.Handle("/articles/view/{slug}", views.ShowArticle{})
	http.Handle("/articles/new/", views.NewArticle{})
	http.Handle("/articles/edit/{slug}", views.EditArticle{})
	http.Handle("/articles/delete/{slug}", http.HandlerFunc(views.DeleteArticle))

	log.Println("serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
