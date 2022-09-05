package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *Application) routes() http.Handler {
	// Initialize the router.
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.clientError(w, http.StatusMethodNotAllowed)
	})

	// Leave the static files route unchanged.
	router.Handler(http.MethodGet, "/static/*filepath", app.NoDirListingHandler(http.Dir(app.StaticDir)))

	// Create a new middleware chain containing the middleware specific to our
	// dynamic application routes. For now, this chain will only contain the
	// LoadAndSave session middleware, but we'll add more to it later.
	dynamic := alice.New(app.SessionManager.LoadAndSave)

	//mux.Handle("/", app.HomeHandler())
	router.Handler(http.MethodGet, "/", dynamic.Then(app.HomeHandler()))

	//mux.Handle("/snippet/view", app.SnippetViewHandler())
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.Then(app.SnippetViewHandler()))

	//mux.Handle("/snippet/create", app.SnippetCreatePostHandler())
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.Then(app.SnippetCreatePostHandler()))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	//standard := alice.New(app.logRequest, secureHeaders)
	// Wrap the router with the middleware and return it as normal.
	return standard.Then(router)
}
