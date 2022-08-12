package main

import (
	"net/http"
)

func routes(app *Application) *http.ServeMux {
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

	mux.Handle("/", HomeHandler(app))
	mux.Handle("/snippet/view", SnippetViewHandler(app))
	mux.Handle("/snippet/create", SnippetCreateHandler(app))
	mux.Handle("/handler/", &handlerImpl{})
	mux.Handle("/snippet/trans", SnippetTransationHandler(app))
	return mux
}