package main

import (
	"github.com/liviu-moraru/snippetbox/internal/models"
	"log"
)

type Application struct {
	InfoLog   *log.Logger
	ErrorLog  *log.Logger
	Snippets  *models.SnippetModel
	StaticDir string
}
