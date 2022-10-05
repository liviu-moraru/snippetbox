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

	// Use the nosurf middleware on all our 'dynamic' routes.
	dynamic := alice.New(app.SessionManager.LoadAndSave, noSurf)

	router.Handler(http.MethodGet, "/", dynamic.Then(app.HomeHandler()))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.Then(app.SnippetViewHandler()))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	// Because the 'protected' middleware chain appends to the 'dynamic' chain
	// the noSurf middleware will also be used on the three routes below too.
	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.Then(app.SnippetCreatePostHandler()))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	//standard := alice.New(app.logRequest, secureHeaders)
	// Wrap the router with the middleware and return it as normal.
	return standard.Then(router)
}
