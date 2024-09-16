package store

import (
	"fmt"
	"net/http"
	"tender-service/models"

	"github.com/jmoiron/sqlx"
)

func ValidateOrganizationId(db *sqlx.DB, organizationID string) (int, error) {
	var valid bool
	err := db.Get(&valid, "SELECT EXISTS (SELECT 1 FROM organization WHERE id = $1)", organizationID)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("ошибка при проверке организации: %v", err)
	}
	if !valid {
		return http.StatusBadRequest, fmt.Errorf("организация с данным ID не существует")
	}
	return http.StatusOK, nil
}

func ValidateUserUsername(db *sqlx.DB, creatorUsername string) (int, error) {
	var valid bool
	err := db.Get(&valid, "SELECT EXISTS (SELECT 1 FROM employee WHERE username = $1)", creatorUsername)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("ошибка при проверке пользователя: %v", err)
	}
	if !valid {
		return http.StatusUnauthorized, fmt.Errorf("пользователь с данным username не существует")
	}
	return http.StatusOK, nil
}

func ValidateUsersAffiliation(db *sqlx.DB, organizationId string, userId string) (int, error) {

	var valid bool
	err := db.Get(&valid, "SELECT EXISTS (SELECT 1 FROM organization_responsible WHERE organization_id = $1 AND user_id = $2)", organizationId, userId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("ошибка при проверке связи пользователя с организацией: %v", err)
	}
	if !valid {
		return http.StatusBadRequest, fmt.Errorf("пользователь не связан с данной организацией")
	}
	return http.StatusOK, nil
}

func ValidUserPermisson(db *sqlx.DB, userId string, tender *models.TenderById) (bool, error) {
	fmt.Println("UserId", userId)
	fmt.Println("tender.Status", tender.Status)
	if tender.UserId == userId || tender.Status == "Published" {
		return true, nil
	}
	_, validErr := ValidateUsersAffiliation(db, tender.OrganizationId, userId)
	if validErr != nil {
		return false, validErr
	}
	return true, nil
}

func ValidTenderPermission(db *sqlx.DB, username string, tenderId string) (*models.TenderById, int, error) {
	userId, userStatus, userErr := GetUserIdByUsername(db, username)
	if userErr != nil {
		return nil, userStatus, fmt.Errorf("ошибка получения пользователя: %v", userErr)
	}
	tender, err := GetTenderById(db, tenderId)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("ошибка получения тендера: %v", err)
	}
	canAccess, err := ValidUserPermisson(db, userId, tender)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if !canAccess {
		return nil, http.StatusUnauthorized, fmt.Errorf("доступ к тендеру запрещен")
	}

	return tender, http.StatusOK, nil
}

func ValidateTenderExist(db *sqlx.DB, tenderId string) (int, error) {
	var valid bool
	err := db.Get(&valid, "SELECT EXISTS (SELECT 1 FROM tender WHERE id = $1)", tenderId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("ошибка при проверке существования тендера: %v", err)
	}
	if !valid {
		return http.StatusBadRequest, fmt.Errorf("тендер не существует")
	}
	return http.StatusOK, nil
}

func ValidateTenderStatus(db *sqlx.DB, tenderId string) (int, error){
	var status string
	err := db.Get(&status, "SELECT status FROM tender WHERE id = $1", tenderId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("ошибка при получении статуса тендера: %v", err)
	}
	if status != "Published" {
		return http.StatusBadRequest, fmt.Errorf("тендер не опубликован")
	}
	return http.StatusOK, nil
}