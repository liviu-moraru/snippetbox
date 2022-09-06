package main

import (
	"github.com/liviu-moraru/snippetbox/internal/models"
	"html/template"
	"path/filepath"
	"time"
)

// Add a Form field with the type "any".
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        any
	Flash       string // Add a flash field to the templateData struct
}

// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time object.
func humanDate(t time.Time) string {
	tz, _ := time.LoadLocation("Local")
	t = t.In(tz)
	return t.Format("02 Jan 06 15:04 -0700")
}

// Initialize a template.FuncMap object and store it in a global variable. This is
// essentially a string-keyed map which acts as a lookup between the names of our
// custom template functions and the functions themselves.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() function to get a slice of all filepaths that
	// match the pattern "./ui/html/pages/*.tmpl". This will essentially gives
	// us a slice of all the filepaths for our application 'page' templates
	// like: [ui/html/pages/home.tmpl ui/html/pages/view.tmpl]
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	// Loop through the page filepaths one-by-one.
	for _, page := range pages {
		// Extract the file name (like 'home.tmpl') from the full filepath
		// and assign it to the name variable.
		name := filepath.Base(page)
		ts, err := parsePage(page)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map, using the name of the page
		// (like 'home.tmpl') as the key.
		cache[name] = ts
	}

	// Return the map.
	return cache, nil
}

func parsePage(page string) (*template.Template, error) {
	name := filepath.Base(page)
	ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
	if err != nil {
		return nil, err
	}

	// Call ParseGlob() *on this template set* to add any partials.
	ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
	if err != nil {
		return nil, err
	}

	// Call ParseFiles() *on this template set* to add the page template.
	ts, err = ts.ParseFiles(page)
	if err != nil {
		return nil, err
	}
	return ts, nil
}
