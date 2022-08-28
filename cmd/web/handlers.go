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
		r.Body = http.MaxBytesReader(w, r.Body, 10)
		var maxBytesError *http.MaxBytesError
		// First we call r.ParseForm() which adds any data in POST request bodies
		// to the r.PostForm map. This also works in the same way for PUT and PATCH
		// requests. If there are any errors, we use our app.ClientError() helper to
		// send a 400 Bad Request response to the user.
		err := r.ParseForm()
		if err != nil {
			app.ErrorLog.Printf("%T", err)
			if errors.As(err, &maxBytesError) {
				app.maxBytesError(w, http.StatusBadRequest)
				return
			}
			app.clientError(w, http.StatusBadRequest)
			return
		}

		// Use the r.PostForm.Get() method to retrieve the title and content
		// from the r.PostForm map.
		title := r.PostForm.Get("title")
		content := r.PostForm.Get("content")

		// The r.PostForm.Get() method always returns the form data as a *string*.
		// However, we're expecting our expires value to be a number, and want to
		// represent it in our Go code as an integer. So we need to manually covert
		// the form data to an integer using strconv.Atoi(), and we send a 400 Bad
		// Request response if the conversion fails.

		expires, err := strconv.Atoi(r.PostForm.Get("expires"))
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
		}

		// Process checkboxes
		for i, item := range r.PostForm["items"] {
			app.InfoLog.Printf("%d: Item %s\n", i, item)
		}

		id, err := app.Snippets.Insert(title, content, expires)
		if err != nil {
			app.serverError(w, err)
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
