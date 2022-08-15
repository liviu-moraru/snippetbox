package main

import (
	"errors"
	"fmt"
	"github.com/liviu-moraru/snippetbox/internal/models"
	"html/template"
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

		files := []string{
			"ui/html/partials/nav.tmpl",
			"ui/html/base.tmpl",
			"ui/html/pages/home.tmpl",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			serverError(app, w, err)
			return
		}

		// Create an instance of a templateData struct holding the snippet data.
		data := &templateData{
			Snippets: snippets,
		}

		err = ts.ExecuteTemplate(w, "base", data)
		if err != nil {
			serverError(app, w, err)
		}

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

		// Initialize a slice containing the paths to the view.tmpl file,
		// plus the base layout and navigation partial that we made earlier.
		files := []string{
			"ui/html/base.tmpl",
			"ui/html/partials/nav.tmpl",
			"ui/html/pages/view.tmpl",
		}

		ts, err := template.ParseFiles(files...)

		// Parse the template files...
		if err != nil {
			serverError(app, w, err)
			return
		}

		/*tz, _ := time.LoadLocation("US/Pacific")
		snippet.Created = snippet.Created.In(tz)*/

		// Create an instance of a templateData struct holding the snippet data.
		data := &templateData{
			Snippet: snippet,
		}
		// And then execute them. Notice how we are passing in the snippet
		// data (a models.Snippet struct) as the final parameter?
		err = ts.ExecuteTemplate(w, "base", data)
		if err != nil {
			serverError(app, w, err)
		}
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
