package main

import (
	"net/http"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(neuteredFileSystem{http.Dir(app.StaticDir)})

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.Handle("/", app.HomeHandler())
	mux.Handle("/snippet/view", app.SnippetViewHandler())
	mux.Handle("/snippet/create", app.SnippetCreateHandler())
	mux.Handle("/handler/", &handlerImpl{})
	mux.Handle("/snippet/trans", app.SnippetTransationHandler())
	return secureHeader(mux)
}
