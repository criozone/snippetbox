package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	dbConn   *sql.DB
}

func main() {
	addr := flag.String("addr", ":8000", "Http network address the server will listen")
	dsn := flag.String("dsn", "root@tcp(localhost:3306)/snippetbox?parseTime=true", "MYSQL data source connection string")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)
	// user:password@tcp(localhost:5555)/dbname?
	dbConn, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatalln(err)
	}

	defer func(dbConn *sql.DB) {
		err := dbConn.Close()
		if err != nil {
			errorLog.Println(err)
		}
	}(dbConn)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		dbConn:   dbConn,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	infoLog.Printf("Server listening in %s\n", *addr)
	err = srv.ListenAndServe()
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
