package domain

import (
	"errors"
	"time"
)

var errNoRecord = errors.New("models: not matching record found")

type Snippet struct {
	Id      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
