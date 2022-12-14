package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"
)

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *Application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Application) customClientError(w http.ResponseWriter, error string, status int) {
	http.Error(w, error, status)
}

func (app *Application) maxBytesError(w http.ResponseWriter, status int) {
	http.Error(w, "Max Bytes Error", status)
}

// For consistency, we'll also implement a notFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found response to
// the user.
func (app *Application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *Application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	var ts *template.Template
	var ok bool
	var err error

	develop := os.Getenv("DEVELOP")

	if strings.ToLower(develop) == "true" {
		fp := filepath.Join("./ui/html/pages", page)
		ts, err = parsePage(fp)
		if err != nil {
			app.serverError(w, fmt.Errorf("the template %s does not exist", page))
			return
		}
	} else {
		ts, ok = app.TemplateCache[page]

		if !ok {
			app.serverError(w, fmt.Errorf("the template %s does not exist", page))
			return
		}
	}

	buf := new(bytes.Buffer)

	// Execute the template set and write the response body. Again, if there
	// is any error we call the serverError() helper.
	err = ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// If the template is written to the buffer without any errors, we are safe
	// to go ahead and write the HTTP status code to http.ResponseWriter.
	w.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWriter. Note: this
	// is another time when we pass our http.ResponseWriter to a function that
	// takes an io.Writer.
	buf.WriteTo(w)
}

// Create an newTemplateData() helper, which returns a pointer to a templateData
// struct initialized with the current year. Note that we're not using the
// *http.Request parameter here at the moment, but we will do later in the book.
func (app *Application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.SessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r), // Add the CSRF token.
	}
}

// Create a new decodePostForm() helper method. The second parameter here, dst,
// is the target destination that we want to decode the form data into.
func (app *Application) decodePostForm(r *http.Request, dst any) error {
	// Call ParseForm() on the request, in the same way that we did in our
	// createSnippetPost handler.
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

// Return true if the current request is from an authenticated user, otherwise
// return false.
func (app *Application) isAuthenticated(r *http.Request) bool {
	return app.SessionManager.Exists(r.Context(), "authenticatedUserID")
}
