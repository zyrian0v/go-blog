package handler

import (
	slugify "github.com/gosimple/slug"
	"go-blog/db"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type FormErrors struct {
	ValidationErrs db.ErrorMap
	DBErr          error
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

	page := 1
	pageQuery := r.FormValue("page")
	if pageQuery != "" {
		var err error
		page, err = strconv.Atoi(pageQuery)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	articles, err := db.GetAllArticles(page)
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
	slug := r.PathValue("slug")
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

type NewArticleView struct {
	FormErrors
}

func (v NewArticleView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		v.post(w, r)
		return
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

	errMap := a.Validate()
	if len(errMap) > 0 {
		v.ValidationErrs = errMap
	} else {
		err := db.AddArticle(a)
		if err != nil {
			v.DBErr = err
		} else {
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
			return
		}
	}

	files := []string{
		"views/layout.html",
		"views/new_article.html",
	}
	tmpl := template.Must(template.ParseFiles(files...))
	if err := tmpl.Execute(w, v); err != nil {
		log.Println(err)
	}
}

type EditArticleView struct {
	FormErrors
	db.Article
}

func (v EditArticleView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		v.post(w, r)
		return
	}

	slug := r.PathValue("slug")

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
	slug := r.PathValue("slug")

	a := db.Article{
		Title:   r.FormValue("title"),
		Slug:    r.FormValue("slug"),
		Content: r.FormValue("content"),
	}
	v.Article = a

	errMap := a.Validate()
	if len(errMap) > 0 {
		v.ValidationErrs = errMap
	} else {
		err := db.EditArticle(slug, a)
		if err != nil {
			v.DBErr = err
		} else {
			http.Redirect(w, r, "/articles/view/"+a.Slug, http.StatusMovedPermanently)
			return
		}
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

func DeleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	slug := r.PathValue("slug")

	log.Printf("delete '%v'", slug)

	if err := db.DeleteArticle(slug); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
