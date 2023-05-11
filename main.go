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

var articles = map[string]db.Article{
	"hello":          {"one", "Hello!", "Hello, World!"},
	"second-article": {"two", "Second article", "Hello, Second!"},
	"third-article": {"three", "Third article!", `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris dignissim eget enim et lobortis. Proin finibus, neque non tincidunt interdum, magna arcu posuere risus, eget consequat urna dolor a nunc. Pellentesque blandit arcu quis suscipit accumsan. Pellentesque ultricies pulvinar commodo. Sed placerat mollis risus, quis mollis libero ornare quis. Donec non ante vitae ipsum aliquet consectetur. Nunc venenatis fringilla consectetur.

Phasellus pharetra malesuada condimentum. Nam quis velit feugiat, efficitur urna finibus, egestas neque. Praesent euismod ligula non ex lobortis lacinia. Morbi viverra ac est vitae semper. Aliquam pulvinar tristique condimentum. Cras lorem ante, suscipit non efficitur et, pharetra nec risus. Proin ac augue urna. Sed a urna quis purus sagittis sollicitudin eget quis metus. Ut posuere eleifend tortor, non iaculis risus bibendum id. Morbi augue metus, tristique vel feugiat at, cursus sit amet leo.

Nulla vitae lobortis sem, id ultricies mi. Etiam condimentum vel odio quis imperdiet. Vestibulum auctor odio eget efficitur sollicitudin. Praesent ut elit purus. Aliquam mollis urna a bibendum vestibulum. Nunc vehicula vehicula tristique. Suspendisse dignissim enim et sagittis molestie. Integer vestibulum magna sed arcu commodo, eget pellentesque sem molestie. Morbi sodales, eros vitae dapibus tempor, diam enim molestie enim, sit amet viverra magna ipsum ut nisi. Nam pharetra sapien justo, ac aliquet est tincidunt a. Fusce fermentum mi non velit aliquam suscipit.

Maecenas nulla lacus, placerat eget commodo nec, dapibus vitae libero. Nunc nec metus quam. Fusce nisi tortor, semper nec cursus sit amet, placerat at dui. Ut sapien tortor, sollicitudin quis sollicitudin non, pulvinar ac urna. Maecenas eget lobortis dui. Donec dignissim ac nisi in elementum. Sed congue nunc urna, ac faucibus velit bibendum in. Nulla interdum pulvinar ipsum, convallis tempus felis ultricies et. Phasellus aliquet bibendum aliquet. Donec maximus elit lobortis, interdum ante et, pharetra sem. Quisque consequat nibh nunc, sit amet scelerisque sem cursus non. Sed tempor diam a diam sollicitudin pellentesque. Praesent tempor, tellus a ornare vulputate, libero libero tempor massa, aliquet varius mi diam ac lectus.

Nullam vel metus at odio eleifend maximus quis non nisl. Praesent efficitur placerat auctor. Vivamus hendrerit lectus vel auctor venenatis. Cras aliquam posuere orci, at fermentum diam tempor eget. Nulla vitae rutrum turpis, ac facilisis diam. Nam consequat magna eu erat vestibulum laoreet. Sed dui nulla, vehicula eu erat sit amet, dictum elementum mauris. Aenean vitae vulputate nibh, suscipit euismod dolor. Sed facilisis maximus eros, at maximus mauris efficitur ut. Phasellus vulputate ultricies odio id tempor. Nunc id congue ante. Donec vitae accumsan diam. In vulputate nisl facilisis quam dapibus gravida.
`},
}

func main() {
	// Database
	db.InitializeHandle()
	db.ApplySchema()

	// Routes
	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	root := IndexView{Intro: "Welcome to my blog!"}
	http.Handle("/", root)

	const articleRoute = "/articles/view/"
	showArticle := http.StripPrefix(articleRoute, ShowArticleView{})
	http.Handle(articleRoute, showArticle)

	const newArticleRoute = "/articles/new/"
	newArticle := http.StripPrefix(newArticleRoute, NewArticleView{})
	http.Handle(newArticleRoute, newArticle)

	const editArticleRoute = "/articles/edit/"
	editArticle := http.StripPrefix(editArticleRoute, EditArticleView{})
	http.Handle(editArticleRoute, editArticle)

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

	files := []string{
		"views/layout.html",
		"views/article.html",
	}
	tmpl := template.Must(template.ParseFiles(files...))

	a, err := db.GetArticleBySlug(slug)
	if err != nil {
		http.NotFound(w, r)
		log.Println(err)
		return
	}
	v.Article = a

	if err := tmpl.Execute(w, v); err != nil {
		log.Println(err)
	}
}

type NewArticleView struct{}

func (v NewArticleView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("new article")

	if r.Method == "POST" {
		a := db.Article{
			Title:   r.FormValue("title"),
			Slug:    slugify.Make(r.FormValue("title")),
			Content: r.FormValue("content"),
		}
		if err := a.Validate(); err != nil {
			fmt.Fprint(w, err)
			return
		}

		err := db.AddArticle(a)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
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
	db.Article
}

func (v EditArticleView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimSuffix(r.URL.Path, "/")
	log.Printf("edit '%v'", slug)

	if r.Method == "POST" {
		a := db.Article{
			Title:   r.FormValue("title"),
			Slug:    r.FormValue("slug"),
			Content: r.FormValue("content"),
		}
		if err := a.Validate(); err != nil {
			fmt.Fprint(w, err)
			return
		}

		err := db.EditArticle(slug, a)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		http.Redirect(w, r, "/articles/view/"+a.Slug, http.StatusMovedPermanently)
	}

	files := []string{
		"views/layout.html",
		"views/edit_article.html",
	}
	tmpl := template.Must(template.ParseFiles(files...))

	a, err := db.GetArticleBySlug(slug)
	if err != nil {
		http.NotFound(w, r)
		log.Println(err)
		return
	}
	v.Article = a

	if err := tmpl.Execute(w, v); err != nil {
		log.Println(err)
	}
}
