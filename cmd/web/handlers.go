package main

import (
	"errors"
	"fmt"
	"github.com/liviu-moraru/snippetbox/internal/models"
	"net/http"
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

func HomeHandler(app *Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			notFound(w)
			return
		}

		snippets, err := app.Snippets.Latest()
		if err != nil {
			serverError(app, w, err)
			return
		}

		// Call the newTemplateData() helper to get a templateData struct containing
		// the 'default' data (which for now is just the current year), and add the
		// snippets slice to it.
		data := newTemplateData(app)
		data.Snippets = snippets

		// Pass the data to the render() helper as normal.
		render(app, w, http.StatusOK, "home.tmpl", data)

	})
}

func SnippetViewHandler(app *Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//testHeaderMap(w)
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			notFound(w)
			return
		}
		snippet, err := app.Snippets.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				notFound(w)
			} else {
				serverError(app, w, err)
			}
			return
		}

		// Call the newTemplateData() helper to get a templateData struct containing
		// the 'default' data (which for now is just the current year), and add the
		// snippets slice to it.
		data := newTemplateData(app)
		data.Snippet = snippet

		render(app, w, http.StatusOK, "view.tmpl", data)

	})
}

// SnippetCreateHandler Add a snippetCreate handler function.
func SnippetCreateHandler(app *Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.Header().Set("Allow", http.MethodPost)
			clientError(w, http.StatusMethodNotAllowed)
			/*w.WriteHeader(405)
			w.Write([]byte("Method not allowed"))*/
			return
		}
		// Create some variables holding dummy data. We'll remove these later on
		// during the build.
		title := "O snail"
		content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
		expires := 7
		// Pass the data to the SnippetModel.Insert() method, receiving the
		// ID of the new record back.
		id, err := app.Snippets.Insert(title, content, expires)
		if err != nil {
			serverError(app, w, err)
		}

		// Redirect the user to the relevant page for the snippet.
		http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	})
}

func SnippetTransationHandler(app *Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.Header().Set("Allow", http.MethodPost)
			clientError(w, http.StatusMethodNotAllowed)
			return
		}

		id, err := app.Snippets.ExampleTransaction()
		if err != nil {
			serverError(app, w, err)
		}

		// Redirect the user to the relevant page for the snippet.
		http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)

	})
}

type handlerImpl struct{}

func (h *handlerImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my handlerImpl"))
}
