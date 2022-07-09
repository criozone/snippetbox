package main

import (
	"criozone.net/snippetbox/pkg/domain"
	"criozone.net/snippetbox/pkg/forms"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippetRep.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tpl", &templateData{Snippets: snippets})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	f := forms.New(r.PostForm)
	f.Required("title", "content", "expires")
	f.MaxLength("title", 100)
	f.Allowed("expires", "365", "7", "1")

	if !f.Valid() {
		app.render(w, r, "create.page.tpl", &templateData{Form: f})

		return
	}

	//TODO: modify form.Form so that we can set filter funcs and call them when getting the field value
	title := strings.TrimSpace(r.PostForm.Get("title"))
	content := strings.TrimSpace(r.PostForm.Get("content"))
	expires := strings.TrimSpace(r.PostForm.Get("expires"))

	id, err := app.snippetRep.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tpl", &templateData{Form: forms.New(nil)})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippetRep.Get(id)
	if err != nil {
		if errors.Is(err, domain.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.tpl", &templateData{Snippet: s})
}

type CustomFileSystem struct {
	fs http.FileSystem
}

func (ncf CustomFileSystem) Open(path string) (http.File, error) {
	f, statErr := ncf.fs.Open(path)
	if statErr != nil {
		return nil, statErr
	}

	stat, statErr := f.Stat()
	if statErr != nil {
		err := f.Close()
		if err != nil {
			return nil, err
		}

		return nil, statErr
	}

	if stat.IsDir() {
		closeErr := f.Close()
		if closeErr != nil {
			return nil, closeErr
		}

		return nil, os.ErrNotExist
	}

	return f, nil
}
