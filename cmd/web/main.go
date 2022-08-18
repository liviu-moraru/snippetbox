package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql" // New import
	"github.com/liviu-moraru/snippetbox/config"
	"github.com/liviu-moraru/snippetbox/internal/models"
	"log"
	"net/http"
	"os"
)

var cfg config.Configuration

func main() {

	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&cfg.DSN, "dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.LUTC|log.Llongfile)

	db, err := openDB(cfg.DSN)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &Application{
		InfoLog:       infoLog,
		ErrorLog:      errorLog,
		Snippets:      &models.SnippetModel{DB: db},
		StaticDir:     cfg.StaticDir,
		TemplateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     cfg.Addr,
		ErrorLog: app.ErrorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s\n", cfg.Addr)
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
