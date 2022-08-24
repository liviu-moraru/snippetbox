package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *Application) routes() http.Handler {
	// Initialize the router.
	router := httprouter.New()
	// We can use PanicHandler instead of a custom middleware( app.recoverPanic)
	/*router.PanicHandler = func(w http.ResponseWriter, r *http.Request, i interface{}) {
		app.serverError(w, fmt.Errorf("%s", i))
	}*/
	// Create a handler function which wraps our notFound() helper, and then
	// assign it as the custom handler for 404 Not Found responses. You can also
	// set a custom handler for 405 Method Not Allowed responses by setting
	// router.MethodNotAllowed in the same way too.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.clientError(w, http.StatusMethodNotAllowed)
	})
	router.Handler(http.MethodGet, "/static/*filepath", app.NoDirListingHandler(http.Dir(app.StaticDir)))

	//mux.Handle("/", app.HomeHandler())
	router.Handler(http.MethodGet, "/", app.HomeHandler())

	//mux.Handle("/snippet/view", app.SnippetViewHandler())
	router.Handler(http.MethodGet, "/snippet/view/:id", app.SnippetViewHandler())

	//mux.Handle("/snippet/create", app.SnippetCreatePostHandler())
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.Handler(http.MethodPost, "/snippet/create", app.SnippetCreatePostHandler())

	// Create the middleware chain as normal.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	//standard := alice.New(app.logRequest, secureHeaders)
	// Wrap the router with the middleware and return it as normal.
	return standard.Then(router)
}
