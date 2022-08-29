package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/liviu-moraru/snippetbox/internal/models"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"unicode/utf8"
)

func (app *Application) HomeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Because httprouter matches the "/" path exactly, we can now remove the
		// manual check of r.URL.Path != "/" from this handler.

		snippets, err := app.Snippets.Latest()
		if err != nil {
			app.serverError(w, err)
			return
		}

		data := app.newTemplateData()
		data.Snippets = snippets

		app.render(w, http.StatusOK, "home.tmpl", data)
	})
}

func (app *Application) SnippetViewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// When httprouter is parsing a request, the values of any named parameters
		// will be stored in the request context. We'll talk about request context
		// in detail later in the book, but for now it's enough to know that you can
		// use the ParamsFromContext() function to retrieve a slice containing these
		// parameter names and values like so:
		params := httprouter.ParamsFromContext(r.Context())

		// We can then use the ByName() method to get the value of the "id" named
		// parameter from the slice and validate it as normal.
		id, err := strconv.Atoi(params.ByName("id"))

		if err != nil || id < 1 {
			app.notFound(w)
			return
		}

		snippet, err := app.Snippets.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
			} else {
				app.serverError(w, err)
			}
			return
		}

		data := app.newTemplateData()
		data.Snippet = snippet

		app.render(w, http.StatusOK, "view.tmpl", data)
	})
}

// Add a new snippetCreate handler, which for now returns a placeholder
// response. We'll update this shortly to show a HTML form.
func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()

	// Initialize a new createSnippetForm instance and pass it to the template.
	// Notice how this is also a great opportunity to set any default or
	// 'initial' values for the form --- here we set the initial value for the
	// snippet expiry to 365 days.
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl", data)
}

// Define a snippetCreateForm struct to represent the form data and validation
// errors for the form fields. Note that all the struct fields are deliberately
// exported (i.e. start with a capital letter). This is because struct fields
// must be exported in order to be read by the html/template package when
// rendering the template.
type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

func (app *Application) SnippetCreatePostHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		// Get the expires value from the form as normal.
		expires, err := strconv.Atoi(r.PostForm.Get("expires"))
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		// Create an instance of the snippetCreateForm struct containing the values
		// from the form and an empty map for any validation errors.
		form := snippetCreateForm{
			Title:       r.PostForm.Get("title"),
			Content:     r.PostForm.Get("content"),
			Expires:     expires,
			FieldErrors: map[string]string{},
		}

		// Update the validation checks so that they operate on the snippetCreateForm
		// instance.
		if strings.TrimSpace(form.Title) == "" {
			form.FieldErrors["title"] = "This field cannot be blank"
		} else if utf8.RuneCountInString(form.Title) > 100 {
			form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
		}

		// Check that the Content value isn't blank.
		if strings.TrimSpace(form.Content) == "" {
			form.FieldErrors["content"] = "This field cannot be blank"
		}

		// Check the expires value matches one of the permitted values (1, 7 or
		// 365).
		if expires != 1 && expires != 7 && expires != 365 {
			form.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
		}

		// If there are any validation errors re-display the create.tmpl template,
		// passing in the snippetCreateForm instance as dynamic data in the Form
		// field. Note that we use the HTTP status code 422 Unprocessable Entity
		// when sending the response to indicate that there was a validation error.
		if len(form.FieldErrors) > 0 {
			data := app.newTemplateData()
			data.Form = &form
			app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
			return
		}

		// We also need to update this line to pass the data from the
		// snippetCreateForm instance to our Insert() method.
		id, err := app.Snippets.Insert(form.Title, form.Content, form.Expires)
		if err != nil {
			app.serverError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	})
}

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
