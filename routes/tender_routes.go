package routes

import (
	"github.com/gorilla/mux"

	"tender-service/app"
	"tender-service/handlers"
)

// NewRouter инициализирует маршруты
func NewTenderRouter(app *app.App, router *mux.Router) {
	router.HandleFunc("", handlers.GetTenders(app)).Methods("GET")
	router.HandleFunc("/new", handlers.CreateTender(app)).Methods("POST")
	router.HandleFunc("/my", handlers.GetUserTenders(app)).Methods("GET")
	router.HandleFunc("/{tenderId}/status", handlers.GetTenderStatus(app)).Methods("GET")
	router.HandleFunc("/{tenderId}/status", handlers.UpdateTenderStatus(app)).Methods("PUT")
	router.HandleFunc("/{tenderId}/edit", handlers.EditTender(app)).Methods("PATCH")
	router.HandleFunc("/{tenderId}/rollback/{version}", handlers.RollbackTender(app)).Methods("PUT")

}
