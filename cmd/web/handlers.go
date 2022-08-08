package main

import (
	"fmt"
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

		err = ts.ExecuteTemplate(w, "base", nil)
		if err != nil {
			serverError(app, w, err)
		}

	})
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	//testHeaderMap(w)
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		notFound(w)
		return
	}
	// Use the fmt.Fprintf() function to interpolate the id value with our response
	// and write it to the http.ResponseWriter.
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// Add a snippetCreate handler function.
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost)
		clientError(w, http.StatusMethodNotAllowed)
		/*w.WriteHeader(405)
		w.Write([]byte("Method not allowed"))*/
	}
	w.Write([]byte("Create a new snippet..."))
}

type handlerImpl struct{}

func (h *handlerImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my handlerImpl"))
}
