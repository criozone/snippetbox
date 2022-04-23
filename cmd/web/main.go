package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	addr := flag.String("addr", ":8000", "Http network address the server will listen")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", home(app))
	mux.HandleFunc("/snippets", listSnippets(app))
	mux.HandleFunc("/snippet", showSnippet(app))
	mux.HandleFunc("/snippet/create", createSnippet(app))
	fileServer := http.FileServer(CustomFileSystem{http.Dir("./ui/static/")})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}
	infoLog.Printf("Server listening in %s\n", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
