package database

import (
    "os"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq" // Импортируем драйвер PostgreSQL
)

// Connect подключается к базе данных
func Connect() (*sqlx.DB, error) {
    connStr := os.Getenv("POSTGRES_CONN")

    db, err := sqlx.Connect("postgres", connStr)
    if err != nil {
        return nil, err
    }
    return db, nil
}