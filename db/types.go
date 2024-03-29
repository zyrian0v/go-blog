package db

import (
	"errors"
	"time"
)

type Article struct {
	Title     string
	Slug      string
	Content   string
	CreatedAt time.Time
}

type ErrorMap map[string][]error

func (a Article) Validate() ErrorMap {
	errs := make(ErrorMap)

	if a.Title == "" {
		err := errors.New("title cant be empty")
		errs["title"] = append(errs["title"], err)
	}
	if a.Slug == "" {
		err := errors.New("slug cant be empty")
		errs["slug"] = append(errs["slug"], err)
	}

	return errs
}
