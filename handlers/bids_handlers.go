package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tender-service/app"
	"tender-service/models"
	"tender-service/store"
	"tender-service/utils"

	"github.com/gorilla/mux"
)

func CreateBid(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request models.BidRequest
		if r.Body == nil {
			utils.WriteErrorResponse(w, "Тело запроса не может быть пустым", http.StatusBadRequest)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			utils.WriteErrorResponse(w, "Некорректное тело запроса", http.StatusBadRequest)
			return
		}


		// Валидация данных на уровне запроса
		if request.Name == "" || request.Description == "" ||
			request.Status == "" || request.TenderId == "" || request.OrganizationId == "" {
			utils.WriteErrorResponse(w, "Отсутствуют обязательные поля", http.StatusBadRequest)
			return
		}

		bid := models.Bid{
			Name:           request.Name,
			Description:    request.Description,
			Status:         "Created",
			Version:        1,
		}


		if (request.OrganizationId != "") {
			if !utils.ValidateAndRespond(w,func() (int, error) { return store.ValidateOrganizationId(app.DB, request.OrganizationId)},) {
				return
			}


			bid.AuthorId = request.OrganizationId
			bid.AuthorType = "Organization"
		}
		if (request.CreatorUsername != "") {

			if !utils.ValidateAndRespond(w, func() (int, error) { return store.ValidateUserUsername(app.DB, request.CreatorUsername) },){
				return
			}

			userId, userStatus, userErr := store.GetUserIdByUsername(app.DB, request.CreatorUsername)
			if userErr != nil {
				utils.WriteErrorResponse(w, userErr.Error(), userStatus)
			}
			if !utils.ValidateAndRespond(w,
			func() (int, error) {
				return store.ValidateUsersAffiliation(app.DB, request.OrganizationId, userId)
			},) {
				return
			}
			bid.AuthorId = userId
			bid.AuthorType = "User"
		} 
		validStatus, validErr:= store.ValidateTenderExist(app.DB, request.TenderId)
		if validErr != nil {
			utils.WriteErrorResponse(w, validErr.Error(), validStatus)
			return
		}
		validStatus, validErr = store.ValidateTenderStatus(app.DB, request.TenderId)
		if validErr != nil {
			utils.WriteErrorResponse(w, validErr.Error(), validStatus)
			return
		} else {
			bid.TenderId = request.TenderId
		}

		newBid, newBidErr := store.CreateBid(app.DB, &bid)
		if newBidErr != nil {
			utils.WriteErrorResponse(w, newBidErr.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(newBid); err != nil {
			utils.WriteErrorResponse(w, "Error encoding response", http.StatusInternalServerError)
		}
	}
}

func GetUserBids(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var limit, offset, paginationStatus int
		var paginationErr error
		limit, offset, paginationStatus, paginationErr = utils.GetPagination(r)
		if paginationErr != nil {
			utils.WriteErrorResponse(w, paginationErr.Error(), paginationStatus)
			return
		}

		username := r.URL.Query().Get("username")

		if username == "" {
			utils.WriteErrorResponse(w, "username is missing", http.StatusBadRequest)
			return
		}

		if !utils.ValidateAndRespond(w,
			func() (int, error) { return store.ValidateUserUsername(app.DB, username) },
		) {
			return
		}

		userId, userStatus, userErr := store.GetUserIdByUsername(app.DB, username)
		if userErr != nil {
			utils.WriteErrorResponse(w, userErr.Error(), userStatus)
			return
		}

		organizationId, organizationIdErr := store.GetOrganizationIdByUserId(app.DB, userId)

		if organizationIdErr != nil {
			utils.WriteErrorResponse(w, organizationIdErr.Error(), http.StatusBadRequest)
			return
		}

		userBids, userBidsErr := store.GetUserBids(app.DB, organizationId, userId, limit, offset)
		if userBidsErr != nil {
			utils.WriteErrorResponse(w, userBidsErr.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("userBids", userBids)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userBids)
	}

}

func GetTenderBids(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var limit, offset, paginationStatus int
		var paginationErr error
		limit, offset, paginationStatus, paginationErr = utils.GetPagination(r)
		if paginationErr != nil {
			utils.WriteErrorResponse(w, paginationErr.Error(), paginationStatus)
			return
		}

		username := r.URL.Query().Get("username")

		if username == "" {
			utils.WriteErrorResponse(w, "username is missing", http.StatusBadRequest)
			return
		}

		vars := mux.Vars(r)
		tenderId := vars["tenderId"]
		if tenderId == "" {
			utils.WriteErrorResponse(w, "Параметры запроса не могут быть пустыми", http.StatusBadRequest)
			return
		}

		validStatus, validErr:= store.ValidateTenderExist(app.DB, tenderId)
		if validErr != nil {
			utils.WriteErrorResponse(w, validErr.Error(), validStatus)
			return
		}

		if !utils.ValidateAndRespond(w,
			func() (int, error) { return store.ValidateUserUsername(app.DB, username) },
		) {
			return
		}
		tenderBids, tenderBidsErr := store.GetTenderBids(app.DB, tenderId, limit, offset)
		if tenderBidsErr != nil {
			utils.WriteErrorResponse(w, tenderBidsErr.Error(), http.StatusInternalServerError)
			return
		}



		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tenderBids)
	}
}



func GetBidStatus(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")

		if username == "" {
			utils.WriteErrorResponse(w, "username is missing", http.StatusBadRequest)
			return
		}

		vars := mux.Vars(r)
		bidId := vars["bidId"]
		if bidId == "" {
			utils.WriteErrorResponse(w, "Параметры запроса не могут быть пустыми", http.StatusBadRequest)
			return
		}
		if !utils.ValidateAndRespond(w,
			func() (int, error) { return store.ValidateUserUsername(app.DB, username) },
		) {
			return
		}
		bidStatus, bidStatusErr := store.GetBidStatus(app.DB, bidId)
		if bidStatusErr != nil {
			utils.WriteErrorResponse(w, bidStatusErr.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bidStatus)
	}
}

func UpdateBidStatus(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")

		if username == "" {
			utils.WriteErrorResponse(w, "username is missing", http.StatusBadRequest)
			return
		}

		status := r.URL.Query().Get("status")

		if status == "" {
			utils.WriteErrorResponse(w, "status is missing", http.StatusBadRequest)
			return
		}

		vars := mux.Vars(r)
		bidId := vars["bidId"]
		if bidId == "" {
			utils.WriteErrorResponse(w, "Параметры запроса не могут быть пустыми", http.StatusBadRequest)
			return
		}
		if !utils.ValidateAndRespond(w,
			func() (int, error) { return store.ValidateUserUsername(app.DB, username) },
		) {
			return
		}

		updatedBid, updatedBidErr := store.UpdateBidStatus(app.DB, bidId, status)
		if updatedBidErr != nil {
			utils.WriteErrorResponse(w, updatedBidErr.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedBid)


	}
}

func EditBid(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bidId := vars["bidId"]
		username := r.URL.Query().Get("username")
		if bidId == "" {
			utils.WriteErrorResponse(w, "Параметры запроса не могут быть пустыми", http.StatusBadRequest)
			return
		}
		if username == ""{
			utils.WriteErrorResponse(w, "Параметр username отсутствует", http.StatusBadRequest)
			return
		}

		if r.Body == nil {
			utils.WriteErrorResponse(w, "Тело запроса не может быть пустым", http.StatusBadRequest)
			return
		}
		// Декодирование тела запроса
		var request models.BidRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			utils.WriteErrorResponse(w, "Некорректное тело запроса", http.StatusBadRequest)
			return
		}
		bidUpdateData := models.Bid{}

		if request.Name != "" {
			bidUpdateData.Name = request.Name
		}
		if request.Description != "" {
			bidUpdateData.Description = request.Description
		}
		updateBid, BidStatus, BidErr := store.UpdateBid(app.DB, bidId, username, &bidUpdateData)
		if BidErr != nil {
			utils.WriteErrorResponse(w, "Ошибка при обновлении предложения: " + BidErr.Error(), BidStatus)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(updateBid); err != nil {
			utils.WriteErrorResponse(w, "Error encoding response", http.StatusInternalServerError)
		}
	}
}

func RollbackBid(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bidId := vars["bidId"]
		version := vars["version"]
		username := r.URL.Query().Get("username")
		if bidId == "" {
			utils.WriteErrorResponse(w, "Параметр id предложения не может быть пустым", http.StatusBadRequest)
			return
		}
		if version == "" || username == "" {
			utils.WriteErrorResponse(w, "Параметр версии не может быть пустым", http.StatusBadRequest)
			return
		}
		fmt.Println("bidId, version", bidId, version)
		rollbackBid, rollbackStatus, rollbackBidErr := store.RollbackBidVersion(app.DB, bidId, version, username)
		if rollbackBidErr != nil {
			utils.WriteErrorResponse(w, "Error updating bid: " + rollbackBidErr.Error(), rollbackStatus)
			return
		}

		// Установка заголовка и отправка ответа
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(rollbackBid); err != nil {
			utils.WriteErrorResponse(w, "Error encoding response", http.StatusInternalServerError)
		}
	}
}