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

type ErrorMap map[string]string

func (a Article) Validate() ErrorMap {
	errs := make(ErrorMap)

	if a.Title == "" {
		errs["title"] = "Title cant be empty"
	}
	if a.Slug == "" {
		errs["slug"] = "Slug cant be empty"
	}

	return errs
}
