package db

import (
	"errors"
)

type Article struct {
	Title   string
	Slug string
	Content string
}


func (a Article) Validate() error {
	if a.Title == "" {
		return errors.New("validate article: title cant be empty")
	}
	if a.Slug == "" {
		return errors.New("validate article: slug cant be empty")
	}

	return nil
}
