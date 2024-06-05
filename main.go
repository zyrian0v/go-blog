package main

import (
	"log"
	"net/http"

	"go-blog/db"
	"go-blog/handler"
)

func main() {
	// Database
	db.InitializeHandle()
	db.ApplySchema()

	// Routes
	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	root := handler.IndexView{Intro: "Welcome to my blog!"}
	http.Handle("/", root)

	http.Handle("/login/", handler.LoginView{})

	// const articleRoute = "/articles/view/"
	// showArticle := http.StripPrefix(articleRoute, handler.ShowArticleView{})
	http.Handle("/articles/view/{slug}", handler.ShowArticleView{})

	// const newArticleRoute = "/articles/new/"
	// newArticle := http.StripPrefix(newArticleRoute, handler.NewArticleView{})
	http.Handle("/articles/new/", handler.NewArticleView{})

	// const editArticleRoute = "/articles/edit/"
	// editArticle := http.StripPrefix(editArticleRoute, handler.EditArticleView{})
	http.Handle("/articles/edit/{slug}", handler.EditArticleView{})

	// const deleteArticleRoute = "/articles/delete/"
	// deleteArticle := http.StripPrefix(deleteArticleRoute, http.HandlerFunc(handler.DeleteArticleHandler))
	http.Handle("/articles/delete/{slug}", http.HandlerFunc(handler.DeleteArticleHandler))

	log.Println("serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
