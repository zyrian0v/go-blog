package main

import (
	"flag"
	"log"
	"net/http"

	"blog/db"
	"blog/views"
)

func main() {
	schemaFlag := flag.Bool("schema", false, "apply schema")
	flag.Parse()

	// Database
	db.InitializeHandle()
	if *schemaFlag {
		db.ApplySchema()
	}

	// Routes
	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.Handle("/", middleware(views.Index{}))

	http.Handle("/login/", middleware(views.Login{}))
	http.Handle("/articles/view/{slug}", middleware(views.ShowArticle{}))
	http.Handle("/articles/new/", middleware(views.NewArticle{}))
	http.Handle("/articles/edit/{slug}", middleware(views.EditArticle{}))
	http.Handle("/articles/delete/{slug}", middleware(http.HandlerFunc(views.DeleteArticle)))

	log.Println("serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func middleware(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}