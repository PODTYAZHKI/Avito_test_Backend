package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"tender-service/app"
	"tender-service/models"
	"tender-service/store"
	"tender-service/utils"

	"github.com/gorilla/mux"
)

// Обработчик получения тендеров
func GetTenders(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("getTenders")
		var limit, offset, paginationStatus int
		var paginationErr error
		limit, offset, paginationStatus, paginationErr = utils.GetPagination(r)
		if paginationErr != nil {
			utils.WriteErrorResponse(w, paginationErr.Error(), paginationStatus)
			return
		}

		// Получаем фильтр по типам услуг
		serviceTypesParam := r.URL.Query()["service_type"]
		var serviceTypes []string
		if len(serviceTypesParam) > 0 {
			serviceTypes = strings.Split(serviceTypesParam[0], ",")
			//TODO
			for _, st := range serviceTypes {
				st = strings.TrimSpace(st)
				if st != "Construction" && st != "Delivery" && st != "Manufacture" {
					utils.WriteErrorResponse(w, fmt.Sprintf("Invalid service_type value: %s", st), http.StatusBadRequest)
					return
				}
			}
		}

		// Извлекаем данные из базы данных
		tenders, err := store.GetTenders(app.DB, serviceTypes, limit, offset)
		fmt.Println("tenders", tenders)

		if err != nil {
			utils.WriteErrorResponse(w, "Error fetching tenders", http.StatusInternalServerError)
			return
		}

		// Возвращаем результат
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tenders)
	}
}

func CreateTender(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request models.TenderRequest

		// Декодирование тела запроса
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			utils.WriteErrorResponse(w, "Некорректное тело запроса", http.StatusBadRequest)
			return
		}

		// Валидация данных на уровне запроса
		if request.Name == "" || request.Description == "" || request.ServiceType == "" ||
			request.Status == "" || request.OrganizationId == "" || request.CreatorUsername == "" {
			utils.WriteErrorResponse(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		// Проверка, существует ли организация и связан ли с ней пользователь
		// Получаем user_id на основе username
		userId, userStatus, userErr := store.GetUserIdByUsername(app.DB, request.CreatorUsername)
		if userErr != nil {
			utils.WriteErrorResponse(w, userErr.Error(), userStatus)
			return
		}

		if !utils.ValidateAndRespond(w,
			func() (int, error) { return store.ValidateOrganizationId(app.DB, request.OrganizationId) },
			func() (int, error) { return store.ValidateUserUsername(app.DB, request.CreatorUsername) },
			func() (int, error) {
				return store.ValidateUsersAffiliation(app.DB, request.OrganizationId, userId)
			},
		) {
			return
		}

		// Создание нового тендера
		tender := models.Tender{
			Name:           request.Name,
			Description:    request.Description,
			ServiceType:    request.ServiceType,
			Status:         "Created",
			OrganizationId: request.OrganizationId,
			UserId:         userId,
			Version:        1,
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
		}

		// Сохранение в базу данных
		newTender, newTenderErr := store.CreateTender(app.DB, &tender)
		if newTenderErr != nil {
			utils.WriteErrorResponse(w, newTenderErr.Error(), http.StatusBadRequest)
			return
		}

		// Установка заголовка и отправка ответа
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(newTender); err != nil {
			utils.WriteErrorResponse(w, "Error encoding response", http.StatusInternalServerError)
		}
	}
}

func GetUserTenders(app *app.App) http.HandlerFunc {
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
		userTenders, userTendersErr := store.GetUserTenders(app.DB, username, limit, offset)
		fmt.Println("User tenders", userTenders)

		if userTendersErr != nil {
			utils.WriteErrorResponse(w, userTendersErr.Error(), http.StatusInternalServerError)
			return
		}

		// Возвращаем результат
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userTenders)

	}
}

func GetTenderStatus(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		tenderId := vars["tenderId"]
		if tenderId == "" {
			utils.WriteErrorResponse(w, "Параметр tenderId отсутствует", http.StatusBadRequest)
			return
		}

		username := r.URL.Query().Get("username")
		if username == "" {
			utils.WriteErrorResponse(w, "Параметр username отсутствует", http.StatusBadRequest)
			return
		}
		if !utils.ValidateAndRespond(w,
			func() (int, error) { return store.ValidateUserUsername(app.DB, username) },
		) {
			return
		}

		tender, tenderStatus, err := store.GetTender(app.DB, tenderId, username)
		if err != nil {
			utils.WriteErrorResponse(w, fmt.Sprintf("Ошибка при получении тендера: %v", err), tenderStatus)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tender.Status)
	}
}

func UpdateTenderStatus(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tenderId := vars["tenderId"]
		status := r.URL.Query().Get("status")
		username := r.URL.Query().Get("username")
		if tenderId == "" {
			utils.WriteErrorResponse(w, "Параметр tenderId отсутствует", http.StatusBadRequest)
			return
		}
		if username == "" || status =="" {
			utils.WriteErrorResponse(w, "Параметр username отсутствует", http.StatusBadRequest)
			return
		}
		_, tenderStatus, err := store.GetTender(app.DB, tenderId, username)
		if err != nil {
			// fmt.Println("TenderStatus", tenderStatus)
			utils.WriteErrorResponse(w, fmt.Sprintf("Ошибка получения тендера: %v", err), tenderStatus)
			return
		}

		// Обновляем статус тендера
		tender, err := store.UpdateTenderStatus(app.DB, tenderId, status)
		if err != nil {
			utils.WriteErrorResponse(w, fmt.Sprintf("Ошибка при обновлении статуса тендера: %v", err), http.StatusInternalServerError)
			return
		}

		// Ответ с успешным обновлением
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tender)
	}
}

func EditTender(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		tenderId := vars["tenderId"]
		username := r.URL.Query().Get("username")
		if tenderId == "" {
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
		var request models.TenderRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			utils.WriteErrorResponse(w, "Некорректное тело запроса", http.StatusBadRequest)
			return
		}
		tenderUpdateData := models.Tender{}

		if request.Name != "" {
			tenderUpdateData.Name = request.Name
		}
		if request.Description != "" {
			tenderUpdateData.Description = request.Description
		}
		if request.ServiceType != "" {
			tenderUpdateData.ServiceType = request.ServiceType
		}

		// Проверка, были ли переданы какие-либо поля для обновления
		// if tenderUpdateData.Name == "" && tenderUpdateData.Description == "" && tenderUpdateData.ServiceType == "" {
		// 	utils.WriteErrorResponse(w, "Не передано никаких параметров для обновления", http.StatusBadRequest)
		// 	return
		// }

		updateTender, tenderStatus, err := store.UpdateTender(app.DB, tenderId, username, &tenderUpdateData)
		if err != nil {
			utils.WriteErrorResponse(w, "Ошибка при обновлении тендера: " + err.Error(), tenderStatus)
			return
		}

		// Установка заголовка и отправка ответа
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(updateTender); err != nil {
			utils.WriteErrorResponse(w, "Error encoding response", http.StatusInternalServerError)
		}

	}
}

func RollbackTender(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tenderId := vars["tenderId"]
		version := vars["version"]
		username := r.URL.Query().Get("username")
		if tenderId == "" {
			utils.WriteErrorResponse(w, "Параметр id тендера не может быть пустым", http.StatusBadRequest)
			return
		}
		if version == "" || username == "" {
			utils.WriteErrorResponse(w, "Параметр версии не может быть пустым", http.StatusBadRequest)
			return
		}
		rollbackTender, rollbackStatus, err := store.RollbackTenderVersion(app.DB, tenderId, version, username)
		if err != nil {
			utils.WriteErrorResponse(w, "Error updating tender: " + err.Error(), rollbackStatus)
			return
		}

		// Установка заголовка и отправка ответа
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(rollbackTender); err != nil {
			utils.WriteErrorResponse(w, "Error encoding response", http.StatusInternalServerError)
		}
	}
}
