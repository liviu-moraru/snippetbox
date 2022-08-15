package main

import (
	"github.com/liviu-moraru/snippetbox/internal/models"
	"html/template"
	"log"
)

type Application struct {
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	Snippets      *models.SnippetModel
	StaticDir     string
	TemplateCache map[string]*template.Template
}
