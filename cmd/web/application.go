package main

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/liviu-moraru/snippetbox/internal/models"
	"html/template"
	"log"
)

// Application Add a new sessionManager field to the application struct.
type Application struct {
	InfoLog        *log.Logger
	ErrorLog       *log.Logger
	Snippets       *models.SnippetModel
	StaticDir      string
	TemplateCache  map[string]*template.Template
	FormDecoder    *form.Decoder
	SessionManager *scs.SessionManager
}
