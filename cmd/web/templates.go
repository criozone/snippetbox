package main

import (
	"criozone.net/snippetbox/pkg/domain"
	"criozone.net/snippetbox/pkg/forms"
	"html/template"
	"path/filepath"
	"time"
)

type templateData struct {
	CurrentYear int
	Form        *forms.Form
	Snippet     *domain.Snippet
	Snippets    []*domain.Snippet
	Flash       string
}

func NewTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func humanDate(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}
