package main

import (
	"github.com/go-playground/form/v4"
	"github.com/liviu-moraru/snippetbox/internal/models"
	"html/template"
	"log"
)

// Application Add a formDecoder field to hold a pointer to a form.Decoder instance.
type Application struct {
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	Snippets      *models.SnippetModel
	StaticDir     string
	TemplateCache map[string]*template.Template
	FormDecoder   *form.Decoder
}
