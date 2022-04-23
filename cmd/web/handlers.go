package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

func home(app *application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		files := []string{
			"./ui/html/home.page.tpl",
			"./ui/html/base.layout.tpl",
			"./ui/html/footer.partial.tpl",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.errorLog.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		err = ts.Execute(w, nil)
		if err != nil {
			app.errorLog.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func createSnippet(app *application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		_, err := w.Write([]byte(fmt.Sprintf("Your are now in %s", r.RequestURI)))
		if err != nil {
			app.errorLog.Println(err)
		}
	}
}

func showSnippet(app *application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			http.NotFound(w, r)
			return
		}

		_, err = fmt.Fprintf(w, "Displaying a snippet with id: %d", id)
		if err != nil {
			app.errorLog.Println(err)
		}
	}
}

func listSnippets(app *application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Snippets list."))
		if err != nil {
			app.errorLog.Println(err)
		}
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
