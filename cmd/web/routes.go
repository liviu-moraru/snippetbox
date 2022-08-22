package main

import (
	"github.com/justinas/alice"
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

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(mux)
	//return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
