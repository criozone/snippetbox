package main

import (
	"criozone.net/snippetbox/pkg/repositories/mysql"
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"github.com/golangcollege/sessions"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	dbConn        *sql.DB
	snippetRep    *mysql.SnippetMysqlRep
	session       *sessions.Session
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":8000", "Http network address the server will listen")
	dsn := flag.String("dsn", "root@tcp(localhost:3306)/snippetbox?parseTime=true", "MYSQL data source connection string")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)
	// user:password@tcp(localhost:5555)/dbname?
	dbConn, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatalln(err)
	}

	defer func(dbConn *sql.DB) {
		fmt.Println("Closing db pool")
		err := dbConn.Close()
		if err != nil {
			errorLog.Println(err)
		}
	}(dbConn)

	tc, err := NewTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatalln(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		dbConn:   dbConn,
		snippetRep: &mysql.SnippetMysqlRep{
			DB: dbConn,
		},
		session:       session,
		templateCache: tc,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	infoLog.Printf("Server listening in %s\n", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
