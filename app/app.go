package app

import (
	"github.com/jmoiron/sqlx"
)

type App struct {
	DB *sqlx.DB
}

// Функции работы с App
func NewApp(db *sqlx.DB) *App {
	return &App{DB: db}
}
