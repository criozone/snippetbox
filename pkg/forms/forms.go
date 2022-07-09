package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

func New(d url.Values) *Form {
	return &Form{
		d,
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		v := strings.TrimSpace(f.Get(field))
		if v == "" {
			f.Errors.Add(field, fmt.Sprintf("Field '%s' should not be empty", field))
		}
	}
}

func (f *Form) MaxLength(field string, d int) {
	v := strings.TrimSpace(f.Get(field))
	if l := utf8.RuneCountInString(v); l > d {
		f.Errors.Add(field, fmt.Sprintf("Value for field '%s' is to long. Max length is %d characters", field, d))
	}
}

func (f *Form) Allowed(field string, opts ...string) {
	v := strings.TrimSpace(f.Get(field))
	for _, opt := range opts {
		if opt == v {
			return
		}
	}

	f.Errors.Add(field, fmt.Sprintf("Field '%s' has invalid content", field))
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
