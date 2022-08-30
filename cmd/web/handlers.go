package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/liviu-moraru/snippetbox/internal/models"
	"github.com/liviu-moraru/snippetbox/internal/validator"
	"net/http"
	"os"
	"path"
	"strconv"
)

func (app *Application) HomeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		params := httprouter.ParamsFromContext(r.Context())

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

func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl", data)
}

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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		// Declare a new empty instance of the snippetCreateForm struct.
		var form snippetCreateForm

		err = app.decodePostForm(r, &form)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		// Then validate and use the data as normal...
		form.CheckField(form.NotBlank(form.Title), "title", "This field cannot be blank")
		form.CheckField(form.MaxCharacters(form.Title, 100), "title", "This field cannot be more than 100 characters long")
		form.CheckField(form.NotBlank(form.Content), "content", "This field cannot be blank")
		form.CheckField(form.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

		// If there are any validation errors re-display the create.tmpl template,
		// passing in the snippetCreateForm instance as dynamic data in the Form
		// field. Note that we use the HTTP status code 422 Unprocessable Entity
		// when sending the response to indicate that there was a validation error.
		if !form.Valid() {
			data := app.newTemplateData()
			data.Form = form
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
