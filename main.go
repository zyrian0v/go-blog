package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"blog-try/db"
	slugify "github.com/gosimple/slug"
)

func main() {
	// Database
	db.InitializeHandle()
	db.ApplySchema()

	// Routes
	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	root := IndexView{Intro: "Welcome to my blog!"}
	http.Handle("/", root)

	http.Handle("/login/", LoginView{})

	const articleRoute = "/articles/view/"
	showArticle := http.StripPrefix(articleRoute, ShowArticleView{})
	http.Handle(articleRoute, showArticle)

	const newArticleRoute = "/articles/new/"
	newArticle := http.StripPrefix(newArticleRoute, NewArticleView{})
	http.Handle(newArticleRoute, newArticle)

	const editArticleRoute = "/articles/edit/"
	editArticle := http.StripPrefix(editArticleRoute, EditArticleView{})
	http.Handle(editArticleRoute, editArticle)

	const deleteArticleRoute = "/articles/delete/"
	deleteArticle := http.StripPrefix(deleteArticleRoute, http.HandlerFunc(deleteArticleHandler))
	http.Handle(deleteArticleRoute, deleteArticle)

	log.Println("serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type IndexView struct {
	Articles []db.Article
	Intro    string
}

func (v IndexView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("root page")

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	articles, err := db.GetAllArticles()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	v.Articles = articles

	files := []string{
		"views/layout.html",
		"views/index.html",
	}
	tmpl := template.Must(template.ParseFiles(files...))
	if err := tmpl.Execute(w, v); err != nil {
		log.Println(err)
		return
	}
}

type ShowArticleView struct {
	db.Article
}

func (v ShowArticleView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimSuffix(r.URL.Path, "/")
	log.Printf("show '%v'", slug)

	a, err := db.GetArticleBySlug(slug)
	if err != nil {
		http.NotFound(w, r)
		log.Println(err)
		return
	}
	v.Article = a

	files := []string{
		"views/layout.html",
		"views/article.html",
	}
	tmpl := template.Must(template.ParseFiles(files...))
	if err := tmpl.Execute(w, v); err != nil {
		log.Println(err)
	}
}

type NewArticleView struct{}

func (v NewArticleView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		v.post(w, r)
	}

	log.Println("new article")

	files := []string{
		"views/layout.html",
		"views/new_article.html",
	}
	tmpl := template.Must(template.ParseFiles(files...))
	if err := tmpl.Execute(w, v); err != nil {
		log.Println(err)
	}
}

func (v NewArticleView) post(w http.ResponseWriter, r *http.Request) {
	a := db.Article{
		Title:   r.FormValue("title"),
		Slug:    slugify.Make(r.FormValue("title")),
		Content: r.FormValue("content"),
	}
	if err := a.Validate(); err != nil {
		fmt.Fprint(w, err)
		return
	}

	if err := db.AddArticle(a); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
	return
}

type EditArticleView struct {
	Errors []error
	db.Article
}

func (v EditArticleView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		v.post(w, r)
		return
	}

	slug := strings.TrimSuffix(r.URL.Path, "/")
	log.Printf("edit '%v'", slug)

	a, err := db.GetArticleBySlug(slug)
	if err != nil {
		http.NotFound(w, r)
		log.Println(err)
		return
	}
	v.Article = a

	files := []string{
		"views/layout.html",
		"views/edit_article.html",
	}
	tmpl := template.Must(template.ParseFiles(files...))
	if err := tmpl.Execute(w, v); err != nil {
		log.Println(err)
	}
}

func (v EditArticleView) post(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimSuffix(r.URL.Path, "/")

	a := db.Article{
		Title:   r.FormValue("title"),
		Slug:    r.FormValue("slug"),
		Content: r.FormValue("content"),
	}
	v.Article = a
	if errs := a.Validate(); errs != nil {
		v.Errors = append(v.Errors, errs...)
	}

	if err := db.EditArticle(slug, a); err != nil {
		v.Errors = append(v.Errors, err)
	}

	if len(v.Errors) == 0 {
		http.Redirect(w, r, "/articles/view/"+a.Slug, http.StatusMovedPermanently)
		return
	}

	files := []string{
		"views/layout.html",
		"views/edit_article.html",
	}
	tmpl := template.Must(template.ParseFiles(files...))
	if err := tmpl.Execute(w, v); err != nil {
		log.Println(err)
	}
}

func deleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	slug := strings.TrimSuffix(r.URL.Path, "/")
	log.Printf("delete '%v'", slug)

	if err := db.DeleteArticle(slug); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

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
