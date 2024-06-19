package views

import (
	"blog/db"
	slugify "github.com/gosimple/slug"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type FormErrors struct {
	ValidationErrs map[string]string
	DBErr          error
}

type Index struct {
	Articles                 []db.Article
	Page, PrevPage, NextPage int
}

func (v Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

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
	v.Page = page
	v.PrevPage = page - 1
	v.NextPage = page + 1

	articles, err := db.GetArticlePage(page)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	v.Articles = articles

	files := []string{
		"templates/layout.html",
		"templates/index.html",
	}
	tmpl := template.Must(template.ParseFiles(files...))
	if err := tmpl.Execute(w, v); err != nil {
		log.Println(err)
		return
	}
}

type ShowArticle struct {
	db.Article
}

func (v ShowArticle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	log.Println(r.URL.Path)

	a, err := db.GetArticleBySlug(slug)
	if err != nil {
		http.NotFound(w, r)
		log.Println(err)
		return
	}
	v.Article = a

	files := []string{
		"templates/layout.html",
		"templates/article.html",
	}
	tmpl := template.Must(template.ParseFiles(files...))
	if err := tmpl.Execute(w, v); err != nil {
		log.Println(err)
	}
}

type NewArticle struct {
	db.Article
	FormErrors
}

func (v NewArticle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

	if r.Method == "POST" {
		v.post(w, r)
		return
	}

	files := []string{
		"templates/layout.html",
		"templates/new_article.html",
	}
	tmpl := template.Must(template.ParseFiles(files...))
	if err := tmpl.Execute(w, v); err != nil {
		log.Println(err)
	}
}

func (v NewArticle) post(w http.ResponseWriter, r *http.Request) {
	a := db.Article{
		Title:   r.FormValue("title"),
		Slug:    slugify.Make(r.FormValue("title")),
		Content: r.FormValue("content"),
	}
	v.Article = a

	files := []string{
		"templates/layout.html",
		"templates/new_article.html",
	}
	tmpl := template.Must(template.ParseFiles(files...))

	errMap := a.Validate()
	if len(errMap) > 0 {
		v.ValidationErrs = errMap
		if err := tmpl.Execute(w, v); err != nil {
			log.Println(err)
		}
		return
	}

	err := db.AddArticle(a)
	if err != nil {
		v.DBErr = err
		if err := tmpl.Execute(w, v); err != nil {
			log.Println(err)
		}
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

type EditArticle struct {
	db.Article
	FormErrors
}

func (v EditArticle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

	if r.Method == "POST" {
		v.post(w, r)
		return
	}

	slug := r.PathValue("slug")
	a, err := db.GetArticleBySlug(slug)
	if err != nil {
		http.NotFound(w, r)
		log.Println(err)
		return
	}
	v.Article = a

	files := []string{
		"templates/layout.html",
		"templates/edit_article.html",
	}
	tmpl := template.Must(template.ParseFiles(files...))
	if err := tmpl.Execute(w, v); err != nil {
		log.Println(err)
	}
}

func (v EditArticle) post(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	a := db.Article{
		Title:   r.FormValue("title"),
		Slug:    r.FormValue("slug"),
		Content: r.FormValue("content"),
	}
	v.Article = a

	files := []string{
		"templates/layout.html",
		"templates/edit_article.html",
	}
	tmpl := template.Must(template.ParseFiles(files...))

	errMap := a.Validate()
	if len(errMap) > 0 {
		v.ValidationErrs = errMap
		if err := tmpl.Execute(w, v); err != nil {
			log.Println(err)
		}
		return
	}

	err := db.EditArticle(slug, a)
	if err != nil {
		v.DBErr = err
		if err := tmpl.Execute(w, v); err != nil {
			log.Println(err)
		}
		return
	}

	http.Redirect(w, r, "/articles/view/"+a.Slug, http.StatusMovedPermanently)
}

func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

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
