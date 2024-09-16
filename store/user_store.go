package store

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func GetUserIdByUsername(db *sqlx.DB, username string) (string, int, error) {
	fmt.Println("user", username)
	var userId string
	err := db.Get(&userId, "SELECT id FROM employee WHERE username = $1", username)
	if err != nil {
		return userId, http.StatusUnauthorized, fmt.Errorf("пользователь с данным username не найден")
	}
	return userId, http.StatusOK, nil
}

