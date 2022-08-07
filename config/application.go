package config

import (
	"database/sql"
	"log"
)

type Application struct {
	InfoLog   *log.Logger
	ErrorLog  *log.Logger
	Addr      string
	StaticDir string
	DSN       string
	DB        *sql.DB
}
