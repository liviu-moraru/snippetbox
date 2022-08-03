package main

import (
	"flag"
	appConfig "github.com/liviu-moraru/snippetbox/config"
	"log"
	"net/http"
	"os"
)

var app appConfig.Application

func main() {

	flag.StringVar(&app.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&app.StaticDir, "static-dir", "./ui/static", "Path to static assets")

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

	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(neuteredFileSystem{http.Dir(app.StaticDir)})

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.Handle("/", HomeHandler(&app))
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	mux.Handle("/handler/", &handlerImpl{})

	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before, and set
	// the ErrorLog field so that the server now uses the custom errorLog logger in
	// the event of any problems.
	srv := &http.Server{
		Addr:     app.Addr,
		ErrorLog: app.ErrorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %s\n", app.Addr)
	// Call the ListenAndServe() method on our new http.Server struct.
	err := srv.ListenAndServe()
	errorLog.Fatal(err)

}
