package middleware

import (
    "log"
    "net/http"
    // "time"
)

// LogRequestMiddleware логирует каждый HTTP-запрос
func LogRequestMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // start := time.Now()

        // Логируем запрос
        log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)

        // Вызов следующего обработчика
        next.ServeHTTP(w, r)

        // Логируем время выполнения
        // log.Printf("Completed in %v", time.Since(start))
    })
}
