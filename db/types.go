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

func (a Article) Validate() (errs []error) {
	if a.Title == "" {
		err := errors.New("title cant be empty")
		errs = append(errs, err)
	}
	if a.Slug == "" {
		err := errors.New("slug cant be empty")
		errs = append(errs, err)
	}

	return
}
