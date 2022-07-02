package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		err := fmt.Errorf("template '%s' not found", name)
		app.serverError(w, err)
		return
	}

	buff := new(bytes.Buffer)
	err := ts.Execute(buff, app.addCommonData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	_, err = buff.WriteTo(w)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) addCommonData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.CurrentYear = time.Now().Year()

	return td
}
