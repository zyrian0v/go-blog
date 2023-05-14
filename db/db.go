package db

import (
	"database/sql"
	"log"
	"os"
	_ "github.com/mattn/go-sqlite3"
)

var handle *sql.DB

func InitializeHandle() {
	var err error
	handle, err = sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	if err := handle.Ping(); err != nil {
		log.Fatal(err)
	}
}

func ApplySchema() {
	schema, err := os.ReadFile("db/schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = handle.Exec(string(schema))
	if err != nil {
		log.Fatal("apply schema: ", err)
	}
}

func GetAllArticles() (as []Article, err error) {
	stmt := `SELECT title, slug, content, created_at 
	FROM articles ORDER BY created_at DESC`
	rows, err := handle.Query(stmt)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		a := Article{}
		err = rows.Scan(&a.Title, &a.Slug, &a.Content, &a.CreatedAt)
		if err != nil {
			return
		}
		as = append(as, a)
	}
	err = rows.Err()

	return
}

func GetArticleBySlug(slug string) (a Article, err error) {
	stmt := `SELECT title, slug, content, created_at
	FROM articles WHERE slug = ?`
	err = handle.QueryRow(stmt, slug).Scan(&a.Title, &a.Slug, &a.Content, &a.CreatedAt)
	return
}

func AddArticle(a Article) (err error) {
	stmt := `INSERT INTO articles (title, slug, content, created_at) 
	VALUES (?, ?, ?, datetime('now', 'localtime'))`
	_, err = handle.Exec(stmt, a.Title, a.Slug, a.Content)
	return
}

func EditArticle(slug string, a Article) (err error) {
	stmt := `UPDATE articles
	SET title = ?,
	slug = ?,
	content = ?
	WHERE slug = ?`
	_, err = handle.Exec(stmt, a.Title, a.Slug, a.Content, slug)
	return
}

func DeleteArticle(slug string) (err error) {
	stmt := "DELETE FROM articles WHERE slug = ?"
	_, err = handle.Exec(stmt, slug)
	return
}
