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

func (app *Application) HomeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			app.notFound(w)
			return
		}

		panic("ooops! sonething went wrong") // Deliberate panic
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
		//testHeaderMap(w)
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
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

func (app *Application) SnippetCreateHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.Header().Set("Allow", http.MethodPost)
			app.clientError(w, http.StatusMethodNotAllowed)
			/*w.WriteHeader(405)
			w.Write([]byte("Method not allowed"))*/
			return
		}

		title := "O snail"
		content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
		expires := 7

		id, err := app.Snippets.Insert(title, content, expires)
		if err != nil {
			app.serverError(w, err)
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	})
}

func (app *Application) SnippetTransationHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.Header().Set("Allow", http.MethodPost)
			app.clientError(w, http.StatusMethodNotAllowed)
			return
		}

		id, err := app.Snippets.ExampleTransaction()
		if err != nil {
			app.serverError(w, err)
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	})
}

type handlerImpl struct{}

func (h *handlerImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my handlerImpl"))
}
