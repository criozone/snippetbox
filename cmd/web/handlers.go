package main

import (
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

	err = ts.Execute(w, nil)
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

	_, err := w.Write([]byte(fmt.Sprintf("Your are now in %s", r.RequestURI)))
	if err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	_, err = fmt.Fprintf(w, "Displaying a snippet with id: %d", id)
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
