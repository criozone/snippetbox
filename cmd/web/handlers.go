package main

import (
	"criozone.net/snippetbox/pkg/domain"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippetRep.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := struct {
		Snippets []*domain.Snippet
	}{
		Snippets: snippets,
	}

	files := []string{
		"./ui/html/home.page.tpl",
		"./ui/html/base.layout.tpl",
		"./ui/html/footer.partial.tpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) listSnippets(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/list.tpl",
		"./ui/html/base.layout.tpl",
		"./ui/html/footer.partial.tpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	dummy := struct {
		List []string
	}{
		List: []string{
			"First dummy snippet",
			"Second dummy snippet",
			"Third dummy snippet",
		},
	}

	err = ts.Execute(w, dummy)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\\nClimb Mount Fuji,\\nBut slowly, slowly!\\n\\n– Kobayashi Issa"
	expires := "7"

	id, err := app.snippetRep.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = fmt.Fprintf(w, "<h2>Displaying a snippet</h2><div><p>%v</p></div>", s)
	if err != nil {
		app.errorLog.Println(err)
	}
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
