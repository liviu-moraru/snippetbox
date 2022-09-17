package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/liviu-moraru/snippetbox/config"
	"github.com/liviu-moraru/snippetbox/internal/models"
	"log"
	"net/http"
	"os"
	"time"
)

var cfg config.Configuration

func main() {

	flag.StringVar(&cfg.Addr, "addr", ":4443", "HTTP network address")
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

	// Initialize a decoder instance...
	formDecoder := form.NewDecoder()

	// Use the scs.New() function to initialize a new session manager. Then we
	// configure it to use our MySQL database as the session store, and set a
	// lifetime of 12 hours (so that sessions automatically expire 12 hours
	// after first being created).
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	/*cookie := &sessionManager.Cookie
	cookie.Name = "mySecondSession"
	cookie.Persist = false*/

	// Make sure that the Secure attribute is set on our session cookies.
	// Setting this means that the cookie will only be sent by a user's web
	// browser when a HTTPS connection is being used (and won't be sent over an
	// unsecure HTTP connection).
	sessionManager.Cookie.Secure = true

	// And add the session manager to our application dependencies.
	app := &Application{
		InfoLog:        infoLog,
		ErrorLog:       errorLog,
		Snippets:       &models.SnippetModel{DB: db},
		StaticDir:      cfg.StaticDir,
		TemplateCache:  templateCache,
		FormDecoder:    formDecoder,
		SessionManager: sessionManager,
	}

	// Initialize a tls.Config struct to hold the non-default TLS settings we
	// want the server to use. In this case the only thing that we're changing
	// is the curve preferences value, so that only elliptic curves with
	// assembly implementations are used.
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// redirect every http request to https
	go http.ListenAndServe(":4000", http.HandlerFunc(httpRedirect))

	// Set the server's TLSConfig field to use the tlsConfig variable we just
	// created.
	srv := &http.Server{
		Addr:      cfg.Addr,
		ErrorLog:  app.ErrorLog,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
		// Add Idle, Read and Write timeouts to the server.
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s\n", cfg.Addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
