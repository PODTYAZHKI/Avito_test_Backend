package routes

import (
	"github.com/gorilla/mux"

	"tender-service/app"
	"tender-service/handlers"
	"tender-service/middleware"
)

func NewDefaultRouter(app *app.App) *mux.Router {
	// Создаем роутер с префиксом /api
	apiRouter := mux.NewRouter().PathPrefix("/api").Subrouter()

	// Эндпоинт для /api/ping
	apiRouter.HandleFunc("/ping", handlers.Ping).Methods("GET")

	// Добавляем подмаршруты для тендеров

	tenders := apiRouter.PathPrefix("/tenders").Subrouter()
	bids := apiRouter.PathPrefix("/bids").Subrouter()

	// Передаем суброутеры для тендеров и предложений
	NewTenderRouter(app, tenders)
	NewBidRouter(app, bids)

	apiRouter.Use(middleware.LogRequestMiddleware)
	return apiRouter
}
