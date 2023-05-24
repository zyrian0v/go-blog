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

	const articleRoute = "/articles/view/"
	showArticle := http.StripPrefix(articleRoute, handler.ShowArticleView{})
	http.Handle(articleRoute, showArticle)

	const newArticleRoute = "/articles/new/"
	newArticle := http.StripPrefix(newArticleRoute, handler.NewArticleView{})
	http.Handle(newArticleRoute, newArticle)

	const editArticleRoute = "/articles/edit/"
	editArticle := http.StripPrefix(editArticleRoute, handler.EditArticleView{})
	http.Handle(editArticleRoute, editArticle)

	const deleteArticleRoute = "/articles/delete/"
	deleteArticle := http.StripPrefix(deleteArticleRoute, http.HandlerFunc(handler.DeleteArticleHandler))
	http.Handle(deleteArticleRoute, deleteArticle)

	log.Println("serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
