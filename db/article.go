package db

import (
	"errors"
	"time"
)

type Timestamp time.Time

func (ts *Timestamp) Scan(v any) error {
	var s string
	switch u := v.(type) {
	case string:
		s = u
	case []byte:
		s = string(u)
	default:
		return errors.New("unsupported type")
	}

	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		return err
	}
	*ts = Timestamp(t)

	return nil
}

func (ts Timestamp) Pretty() string {
	t := time.Time(ts)
	s := t.Format("01.02.2006")
	return s
}

type Article struct {
	Title     string
	Slug      string
	Content   string
	CreatedAt Timestamp
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
