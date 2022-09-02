# Chapter 7. Advanced routing

# 7.2 Clean URLs and method-based routing

- It was used `github.com/julienschmidt/httprouter`  package.

```
go get github.com/julienschmidt/httprouter
```

- Retrieve parameters from context:

```
params := httprouter.ParamsFromContext(r.Context())
id, err := strconv.Atoi(params.ByName("id"))
```

- I changed the way to implement the no directory listing

```
cmd/web/handlers.go

func (app *Application) NoDirListingHandler(d http.Dir) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		fp := params.ByName("filepath")
		fp = path.Join(string(d), fp)
		if fi, err := os.Stat(fp); err == nil && !fi.IsDir() {
			fileServer := http.FileServer(d)
			http.StripPrefix("/static", fileServer).ServeHTTP(w, r)
			return
		}
		app.notFound(w)
	})
}

cmd/web/routes.go

router.Handler(http.MethodGet, "/static/*filepath", app.NoDirListingHandler(http.Dir(app.StaticDir)))

```

- Create a handler function as a custom handler for 404 Not Found or other responses

```
router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
```

- Using query parameters with httprouter:

The router matches against the request methjod and request URL only. You can get query parameter values:

```
handlers.go

func (app *Application) SnippetViewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

}

routes.go

router.Handler(http.MethodGet, "/snippet/view", app.SnippetViewHandler())

```

- Using httprouter.PanicHandler

We can get rid of custom app.recoverPanic middleware:

```
route.go
...
//standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
standard := alice.New(app.logRequest, secureHeaders)
	// Wrap the router with the middleware and return it as normal.
return standard.Then(router)
```

and instead set PanicHandler:

```
routes.go

router := httprouter.New()
router.PanicHandler = func(w http.ResponseWriter, r *http.Request, i interface{}) {
		app.serverError(w, fmt.Errorf("%s", i))
}
```