package forms

import (
	"fmt"
	"net/url"
	"unicode/utf8"
)

type Form struct {
	ParsedForm  url.Values
	Errors      errors
	formFilters []FilterFunc
}

func New(d url.Values) *Form {
	return &Form{
		d,
		errors(map[string][]string{}),
		[]FilterFunc{},
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		v := f.Get(field)
		if v == "" {
			f.Errors.Add(field, fmt.Sprintf("Field '%s' should not be empty", field))
		}
	}
}

func (f *Form) MaxLength(field string, d int) {
	v := f.Get(field)
	if l := utf8.RuneCountInString(v); l > d {
		f.Errors.Add(field, fmt.Sprintf("Value for field '%s' is to long. Max length is %d characters", field, d))
	}
}

func (f *Form) Allowed(field string, opts ...string) {
	v := f.Get(field)
	for _, opt := range opts {
		if opt == v {
			return
		}
	}

	f.Errors.Add(field, fmt.Sprintf("Field '%s' has invalid content", field))
}

func (f *Form) Get(field string) string {
	return f.filter(f.ParsedForm.Get(field))
}

func (f *Form) AddFilter(ff FilterFunc) {
	f.formFilters = append(f.formFilters, ff)
}

func (f *Form) filter(s string) string {
	for _, c := range f.formFilters {
		s = c(s)
	}

	return s
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
