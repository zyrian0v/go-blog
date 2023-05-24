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
	s := t.Format("02.01.2006")
	return s
}

type Article struct {
	Title     string
	Slug      string
	Content   string
	CreatedAt Timestamp
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
