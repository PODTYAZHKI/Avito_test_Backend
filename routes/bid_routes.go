package routes

import (
	"github.com/gorilla/mux"

	"tender-service/app"
	"tender-service/handlers"
)

func NewBidRouter(app *app.App, router *mux.Router) {

	router.HandleFunc("/new", handlers.CreateBid(app)).Methods("POST")
	router.HandleFunc("/my", handlers.GetUserBids(app)).Methods("GET")
	router.HandleFunc("/{tenderId}/list", handlers.GetTenderBids(app)).Methods("GET")
	router.HandleFunc("/{bidId}/status", handlers.GetBidStatus(app)).Methods("GET")
	router.HandleFunc("/{bidId}/status", handlers.UpdateBidStatus(app)).Methods("PUT")
	router.HandleFunc("/{bidId}/edit", handlers.EditBid(app)).Methods("PATCH")
	router.HandleFunc("/{bidId}/rollback/{version}", handlers.RollbackBid(app)).Methods("PUT")

}
