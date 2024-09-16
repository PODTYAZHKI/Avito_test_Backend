package handlers

import (
    "net/http"
)

// Обработчик пинга
func Ping(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Ok"))
}