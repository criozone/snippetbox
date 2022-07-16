package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	md := alice.New(app.recoverPanic, app.requestLogger, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))
	mux.Get("/snippet/delete/:id", dynamicMiddleware.ThenFunc(app.delSnippet))

	fileServer := http.FileServer(CustomFileSystem{http.Dir("./ui/static/")})
	//mux.Get("/static", http.NotFoundHandler())
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return md.Then(mux)
}
