package db

import (
	"time"
)

type Article struct {
	Title     string
	Slug      string
	Content   string
	CreatedAt time.Time
}

func (a Article) Validate() map[string]string {
	errs := make(map[string]string)

	if a.Title == "" {
		errs["title"] = "Title cant be empty"
	}
	if a.Slug == "" {
		errs["slug"] = "Slug cant be empty"
	}

	return errs
}
