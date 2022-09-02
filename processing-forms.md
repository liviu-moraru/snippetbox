# Chapter 8. Processing forms

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