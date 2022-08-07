package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql" // New import
	appConfig "github.com/liviu-moraru/snippetbox/config"
	"log"
	"net/http"
	"os"
)

var app appConfig.Application

func main() {

	flag.StringVar(&app.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&app.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	// Define a new command-line flag for the MySQL DSN string.
	flag.StringVar(&app.DSN, "dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	//addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// Use log.New() to create a logger for writing information messages. This takes
	// three parameters: the destination to write the logs to (os.Stdout), a string
	// prefix for message (INFO followed by a tab), and flags to indicate what
	// additional information to include (local dat
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC)

	// Create a logger for writing error messages in the same way, but use stderr as
	// the destination and use the log.Lshortfile flag to include the relevant
	// file name and line number.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.LUTC|log.Llongfile)
	app.ErrorLog = errorLog
	app.InfoLog = infoLog

	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate openDB() function below. We pass openDB() the DSN
	// from the command-line flag.
	db, err := openDB(app.DSN)
	if err != nil {
		errorLog.Fatal(err)
	}

	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.

	defer db.Close()
	app.DB = db

	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before, and set
	// the ErrorLog field so that the server now uses the custom errorLog logger in
	// the event of any problems.
	srv := &http.Server{
		Addr:     app.Addr,
		ErrorLog: app.ErrorLog,
		Handler:  routes(&app),
	}

	infoLog.Printf("Starting server on %s\n", app.Addr)
	// Call the ListenAndServe() method on our new http.Server struct.
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
