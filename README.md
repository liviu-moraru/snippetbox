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

# 8.1 Setting up a HTML form

- I modified code in order to not use template caching when DEVELOP environment variable is set with the value true.
  It's better for development. The application must not be started if a change is made in a template file.

# 8.2 Parsing form data

- First we call http.Request.ParseForm() which adds any data in POST request bodies to the r.PostForm map. This also works in the same way for PUT and PATCH requests.

```
func (app *Application) SnippetCreatePostHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
```

- In our code above, we accessed the form values via the map. But an alternative r.PostForm approach is to use the (subtly different) r.Form map
- The r.PostForm map is populated only for POST, PATCH and PUT requests, and contains the
  form data from the request body.
- In contrast, the r.Form map is populated for all requests (irrespective of their HTTP method),
  and contains the form data from any request body and any query string parameters. So, if our
  form was submitted to /snippet/create?foo=bar, we could also get the value of the foo
  parameter by calling r.Form.Get("foo")
- Strictly speaking, the r.PostForm.Get() method that we’ve used above only returns the first
  value for a specific form field. This means you can’t use it with form fields which potentially
  send multiple values, such as a group of checkboxes. Sample of using checkboxes:

```
create.tmpl

 <!-- Checkboxes -->
        <div>
            <input type="checkbox" name="items" value="foo"> Foo
            <input type="checkbox" name="items" value="bar"> Bar
            <input type="checkbox" name="items" value="baz"> Baz
        </div>
  
 handlers.go
 
		for i, item := range r.PostForm["items"] {
			app.InfoLog.Printf("%d: Item %s\n", i, item)
		}
```

- Limiting form size:

  - Unless you’re sending multipart data (i.e. your form has the enctype="multipart/form-data"
    attribute) then POST, PUT and PATCH request bodies are limited to 10MB. If this is exceeded
    then r.ParseForm() will return an error.
  - If you want to change this limit you can use the http.MaxBytesReader() function like so:

```
// Limit the request body size to 4096 bytes
r.Body = http.MaxBytesReader(w, r.Body, 4096)
err := r.ParseForm()
if err != nil {
   http.Error(w, "Bad Request", http.StatusBadRequest)
   return
}
```

- With this code only the first 4096 bytes of the request body will be read during
  r.ParseForm(). Trying to read beyond this limit will cause the MaxBytesReader to return an
  error, which will subsequently be surfaced by r.ParseForm().
  Additionally — if the limit is hit — MaxBytesReader sets a flag on http.ResponseWriter which
  instructs the server to close the underlying TCP connection.
- For Go versions <1.19, the error returned by r.ParseForm is of the errors.errorString type and can be checked only by the text: **http: request body too large**
- In Go 1.19 the struct http.MaxBytesReader was introduces. So, to check this error:

```
handlers.go

func (app *Application) SnippetCreatePostHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 10)
		var maxBytesError *http.MaxBytesError
		err := r.ParseForm()
		if err != nil {
			if errors.As(err, &maxBytesError) {
				app.maxBytesError(w, http.StatusBadRequest)
				return
			}
			app.clientError(w, http.StatusBadRequest)
			return
		}

helpers.go

func (app *Application) maxBytesError(w http.ResponseWriter, status int) {
	http.Error(w, "Max Bytes Error", status)
}
```

# 8.3 Validating form data

- When we check the length of the title field, we’re using the
  utf8.RuneCountInString() function — not Go’s len() function
- Patterns for processing and validating different types of inputs: [this blog post](https://www.alexedwards.net/blog/validation-snippets-for-go)

# 8.4 Displaying errors and repopulating fields

- For the validation errors, the underlying type of our FieldErrors field is a
  map[string]string, which uses the form field names as keys. **For maps, it’s possible to
  access the value for a given key by simply chaining the key name**. So, for example, to render a
  validation error for the title field we can use the tag {{.Form.FieldErrors.title}} in our
  template.
- Note: Unlike struct fields, map key names don’t have to be capitalized in order to access
  them from a template.
- All the snippetCreateForm struct fields are deliberately exported (i.e. start with a capital letter). **This is because struct fields must be exported in order to be read by the html/template package when rendering the template**.
- In handlers.go, line 146 we can pass the address of the form variable instead of a copy of the struct.

```
if len(form.FieldErrors) > 0 {
data := app.newTemplateData()
data.Form = &form // Instead of data.Form = form
app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
return
...
}
```

# 8.5 Creating validation helpers

- Embedding structs in Go, see [good introduction](https://eli.thegreenplace.net/2020/embedding-in-go-part-1-structs-in-structs/)

# 8.6 Automating form parsing

- Install package go-playground/form (see other package gorilla/schema)

```
go get github.com/go-playground/form/v4@v4
```
- Add a pointer to Decoder as dependency

```
type Application struct {
	....
	FormDecoder   *form.Decoder
}
```
- Update snippetCreateForm struct to include struct tags which tell the decoder how to map HTML form values into the different struct fields.

```
handlers.go

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires""`
	validator.Validator `form:"-"`
}
```

- Add a helper function to decode form
```
func (app *Application) decodePostForm(r *http.Request, dst any) error {
    err := r.ParseForm()
	if err != nil {
		return err
	}
	// Call Decode() on our decoder instance, passing the target destination as
	// the first parameter.
	err = app.FormDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// If we try to use an invalid target destination, the Decode() method
		// will return an error with the type *form.InvalidDecoderError.We use
		// errors.As() to check for this and raise a panic rather than returning
		// the error.
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		// For all other errors, we return them as normal.
		return err
	}
	return nil
}
```
Note: Using the erros.As: the form.InvalidDecoderError implements error interface as a method pointer receiver. So, the correct way to test if error is from the specific type is:

```
var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
```
If it had been declared as such:
```
var invalidDecoderError form.InvalidDecoderError
```
the errors.As would raise a panic, because form.InvalidDecoderError doesn't implement error as a value receiver.

- Insert the Decoder dependency

```
main.go

formDecoder := form.NewDecoder()

	// And add it to the application dependencies.
	app := &Application{
		InfoLog:       infoLog,
		ErrorLog:      errorLog,
		Snippets:      &models.SnippetModel{DB: db},
		StaticDir:     cfg.StaticDir,
		TemplateCache: templateCache,
		FormDecoder:   formDecoder,
	}

```
- Using the helper method in handlers:

```
handlers.go

// Update our snippetCreateForm struct to include struct tags which tell the
// decoder how to map HTML form values into the different struct fields. So, for
// example, here we're telling the decoder to store the value from the HTML form
// input with the name "title" in the Title field. The struct tag `form:"-"`
// tells the decoder to completely ignore a field during decoding.
type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires""`
	validator.Validator `form:"-"`
}

func (app *Application) SnippetCreatePostHandler() http.Handler {
  ...
  // Declare a new empty instance of the snippetCreateForm struct.
  var form snippetCreateForm
  
  err = app.decodePostForm(r, &form)
  if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
   }
   .....
  
}
```