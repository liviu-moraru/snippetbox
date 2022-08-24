# 4.2 Installing a database driver

```shell
go get github.com/go-sql-driver/mysql@v1
# Remove package ( from module and from installed packages $GOPATH/pkg/mod)
go get github.com/go-sql-driver/mysql@none
```

# 4.3 Modules and reproducible builds

```shell
# Verify if the checksum s of the downloaded packages mathch the entries in go.sum
go mod verify
# Download the exact versions of all the packages in the project
go mod download
# Upgrade packages to the latest version
go get -u github.com/foo/bar
# Or alternatively, if you want to upgrade to a specific version
go get -u github.com/foo/bar@v2.0.0
# Removing unused packages
go get github.com/foo/bar@none
# go mod tidy will automatically remove any unused packages from your go.mod and go.sum files
# go mod tidy doesn't remove the modules from $GOPATH/bin/mod
go mod tidy -v # -v flag causes tidy to print information about removed modules to standard error.
```

# 4.6 Executing SQL statements

To test create request:

```shell
curl -iL -X POST http://localhost:4000/snippet/create

# Test inside container
# With container mysql running 
docker run -e MYSQL_ROOT_PASSWORD=my-passw --name mysql -p 3306:3306 -v mysql:/var/lib/mysql mysql
( or docker start mysql)

docker exec -it mysql mysql -uroot -p
# Insert password for root (my-passw)
# Inside mysql client REPL
use snippetbox; select * from snippets;
```

To test view request:

```shell
docker start mysql
curl -iL "http://127.0.0.1:4000/snippet/view?id=1"
```

# 4.9 Transactions and other details. Transaction and DB tests.

```shell
go test -v ./internal/models
# or
go test -v ./...
```

# 5.1 Displaying dynamic data

1. DB NULL values in templates

```
<strong>{{.Title.Value}}</strong>
or
 <strong>{{if .Newcol.Valid}}{{.Title.String}}{{else}}-{{end}}</strong>

```

2. Time fields

[Time package](https://pkg.go.dev/time#LoadLocation)

Create an other Time for a time zone

```
t := time.Now()
tz, _ := time.LoadLocation("America/Toronto")
t = t.In(tz)
```

Ex. of time zone: See the file `$GOROOT/lib/time/zoneinfo.zip`

**Format date/time**
See [Format](https://pkg.go.dev/time#Time.Format)

In `format.go` (package time), see the constants:

```
const (
Layout      = "01/02 03:04:05PM '06 -0700" // The reference time, in numerical order.
ANSIC       = "Mon Jan _2 15:04:05 2006"
UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
RFC822      = "02 Jan 06 15:04 MST"
RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
RFC3339     = "2006-01-02T15:04:05Z07:00"
RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
Kitchen     = "3:04PM"
// Handy time stamps.
Stamp      = "Jan _2 15:04:05"
StampMilli = "Jan _2 15:04:05.000"
StampMicro = "Jan _2 15:04:05.000000"
StampNano  = "Jan _2 15:04:05.000000000"
)
```

Ex in templates:

```
 <time>Created: {{.Created.Format "02 Jan 06 15:04 -0700"}}</time>
```

3. Jetbrains Goland: Associate *.tmpl files with Go templates

Preferences -> Editor -> File Types

Then for Go template, add *.tmpl in the list of associated file types.

# 5.2 Template actions and functions

List of template functions: [Template functions](https://pkg.go.dev/text/template#hdr-Functions)

# 6.3 Request logging

- How to log response code

See: [Logging the status code of a HTTP Handler in Go](https://dev.to/julienp/logging-the-status-code-of-a-http-handler-in-go-25aa)

```
// cmd/web/middleware.go

...
type StatusRecorder struct {
    http.ResponseWriter
    Status int
}

//Override the WriteHeader method of the embedded ResponseWriter
func (sr *StatusRecorder) WriteHeader(status int) {
    sr.Status = status
  
    // Without this, the Status Code of the response would not be set.
    sr.ResponseWriter.WriteHeader(status)
}

func (app *Application) logRequest(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
sr := &StatusRecorder{
ResponseWriter: w,
Status:         200,
}

rd := fmt.Sprintf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
next.ServeHTTP(sr, r) // the sr.Status will be set by the WriteHeader method
app.InfoLog.Printf("%s Response status: %d", rd, sr.Status)
})
}
```

```
// cmd/web/routes..go

func (app *Application) routes() http.Handler {
    mux := http.NewServeMux()
	....
    return app.logRequest(secureHeader(mux))
}
```

```
// cnd/web/main.go
...
srv := &http.Server{
    Addr:     cfg.Addr,
    ErrorLog: app.ErrorLog,
    Handler:  app.routes(),
}

err = srv.ListenAndServe()
errorLog.Fatal(err)

```

# 6.4 Panic recovery

1. Setting the **Connection: Close** header on the response acts as a trigger to make Go’ HTTP server automatically close the current connection after a response has been sent.I talso informs the user that the connection will be closed. Note: If the protocol being used is HTTP/2, Go will automatically strip the **Connection: Close** header from the response (so it is not malformed) and send a GOAWAY frame.
2. Go’s HTTP server assumes that the effect of any panic is isolated to the goroutine serving the active HTTP request. Specifically, following a panic our server will log a stack trace to the server error log, unwind
   the stack for the affected goroutine (calling any deferred functions along the way) and close
   the underlying HTTP connection.
3. A panic in a function called as a goroutine by a handler, will end the application (web server will crash)

```
func panicTest() {
	panic("Panic")
}

func (app *Application) SnippetViewHandler() http.Handler { 
...
go panicTest()
}
```

# 6.5 Composable middleware chains

Install alice package:

```
go get github.com/justinas/alice@v1
```

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
