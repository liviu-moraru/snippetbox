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

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(name string) (file http.File, err error) {
	f, err := nfs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		index := path.Join(name, "index.html")
		indexFile, err := nfs.fs.Open(index)
		if err != nil {
			if closeError := f.Close(); closeError != nil {
				return nil, closeError
			}
			return nil, err
		}

		if err := indexFile.Close(); err != nil {
			return nil, err
		}
	}
	return f, nil

}

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

	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *Application) SnippetCreatePostHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		title := r.PostForm.Get("title")
		content := r.PostForm.Get("content")

		expires, err := strconv.Atoi(r.PostForm.Get("expires"))
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			// Correct error
			return
		}

		// Initialize a map to hold any validation errors for the form fields.
		fieldErrors := make(map[string]string)

		// Check that the title value is not blank and is not more than 100
		// characters long. If it fails either of those checks, add a message to the
		// errors map using the field name as the key.
		if strings.TrimSpace(title) == "" {
			fieldErrors["title"] = "This field cannot be blank"
		} else if utf8.RuneCountInString(title) > 100 {
			fieldErrors["title"] = "This field cannot be more than 100 characters long"
		}

		// Check that the Content value isn't blank.
		if strings.TrimSpace(content) == "" {
			fieldErrors["content"] = "This field cannot be blank"
		}

		// Check the expires value matches one of the permitted values (1, 7 or
		// 365).
		if expires != 1 && expires != 7 && expires != 365 {
			fieldErrors["content"] = "This field must equal 1, 7 or 365"
		}

		// If there are any errors, dump them in a plain text HTTP response and
		// return from the handler.

		if len(fieldErrors) > 0 {
			app.customClientError(w, fmt.Sprintf("%v", fieldErrors), http.StatusBadRequest)
			return
		}

		id, err := app.Snippets.Insert(title, content, expires)
		if err != nil {
			app.serverError(w, err)
			// Correct error
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
