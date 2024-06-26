package views

import (
	"blog/db"
	"context"
	"fmt"
	"github.com/go-session/session/v3"
	slugify "github.com/gosimple/slug"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

const (
	paginateBy = 5
)

type FormErrors struct {
	ValidationErrs map[string]string
	DBErr          error
}

type Auth struct {
	User string
}

type Index struct {
	Auth
	Flash                               string
	Articles                            []db.Article
	Page, PrevPage, NextPage, PageCount int
}

func (v Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	flash, ok := store.Get("flash")
	if ok {
		v.Flash = flash.(string)
		store.Delete("flash")
		err := store.Save()
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	}
	user, ok := store.Get("user")
	if ok {
		v.User = user.(string)
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
	count, err := db.GetArticleCount()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	pageCount := count / paginateBy
	if count%paginateBy != 0 {
		pageCount += 1
	}
	v.PageCount = pageCount

	articles, err := db.GetArticlePage(page, paginateBy)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	v.Articles = articles

	files := []string{
		"templates/base.html",
		"templates/index.html",
	}
	tmpl := template.Must(template.ParseFiles(files...))
	if err := tmpl.Execute(w, v); err != nil {
		log.Println(err)
		return
	}
}

type ShowArticle struct {
	Auth
	db.Article
}

func (v ShowArticle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	user, ok := store.Get("user")
	if ok {
		v.User = user.(string)
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
		"templates/base.html",
		"templates/article.html",
	}
	tmpl := template.Must(template.ParseFiles(files...))
	if err := tmpl.Execute(w, v); err != nil {
		log.Println(err)
	}
}

type NewArticle struct {
	Auth
	db.Article
	FormErrors
}

func (v NewArticle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	user, ok := store.Get("user")
	if ok {
		v.User = user.(string)
	}

	if v.User == "" {
		http.Error(w, "Not authenticated", 500)
		return
	}

	if r.Method == "POST" {
		v.post(w, r)
		return
	}

	files := []string{
		"templates/base.html",
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
		"templates/base.html",
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

	http.Redirect(w, r, "/articles/view/"+a.Slug, 303)
}

type EditArticle struct {
	Auth
	db.Article
	FormErrors
}

func (v EditArticle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	user, ok := store.Get("user")
	if ok {
		v.User = user.(string)
	}

	if v.User == "" {
		http.Error(w, "Not authenticated", 500)
		return
	}

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
		"templates/base.html",
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
		"templates/base.html",
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

	http.Redirect(w, r, "/articles/view/"+a.Slug, 303)
}

type DeleteArticle struct {
	Auth
}

func (v DeleteArticle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	user, ok := store.Get("user")
	if ok {
		v.User = user.(string)
	}

	if v.User == "" {
		http.Error(w, "Not authenticated", 500)
		return
	}

	slug := r.PathValue("slug")
	if err := db.DeleteArticle(slug); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	flash := fmt.Sprintf("Successfully deleted \"%v\".", slug)
	store.Set("flash", flash)
	http.Redirect(w, r, "/", 303)
}
