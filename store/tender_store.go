package store

import (
	// "fmt"
	"fmt"
	"net/http"
	// "net/http"
	"tender-service/models"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func GetTenders(db *sqlx.DB, serviceTypes []string, limit int, offset int) ([]models.Tender, error) {
	var tenders []models.Tender
	query := `
	SELECT id, name, description, status, service_type, version, created_at
	FROM tender
	WHERE status = 'Published'`
	args := []interface{}{limit, offset} // Аргументы по умолчанию

	// Фильтр по типу услуг
	if len(serviceTypes) > 0 {
		query += " AND service_type = ANY($1) ORDER BY name ASC LIMIT $2 OFFSET $3"
		args = append([]interface{}{pq.StringArray(serviceTypes)}, args...) // Добавляем serviceTypes в начало
	} else {
		query += " ORDER BY name ASC LIMIT $1 OFFSET $2"
	}
	fmt.Println("Query:", query)
	fmt.Println("Args:", args)
	// Выполнение запроса
	err := db.Select(&tenders, query, args...)
	if err != nil {
		return nil, err
	}
	if len(tenders) == 0 {
		// Если нет результатов, вернуть пустой массив
		return []models.Tender{}, nil
	}

	return tenders, nil
}

func CreateTender(db *sqlx.DB, tender *models.Tender) (*models.Tender, error) {
	query := `
		INSERT INTO tender (name, description, service_type, status, version, organization_id, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at`

	err := db.QueryRowx(query, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.Version, tender.OrganizationId, tender.UserId, tender.CreatedAt, tender.UpdatedAt).StructScan(tender)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании тендеров: %v", err)
	}

	return tender, nil
}

func GetUserTenders(db *sqlx.DB, username string, limit int, offset int) (*[]models.Tender, error) {
	userId, _, userErr := GetUserIdByUsername(db, username)
	if userErr != nil {
		return nil, userErr
	}
	query := `
        SELECT id, name, description, status, service_type, version, created_at
        FROM tender
        WHERE user_id = $1
        ORDER BY name ASC LIMIT $2 OFFSET $3
    `
	var tenders []models.Tender
	// Выполнение запроса
	err := db.Select(&tenders, query, userId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса на получение тендеров: %v", err)
	}

	return &tenders, nil

}

func GetTenderById(db *sqlx.DB, tenderId string) (*models.TenderById, error) {
	query := `
        SELECT id, name, description, status, organization_id, user_id, service_type, version, created_at, updated_at
        FROM tender
        WHERE id = $1
    `
	var tender models.TenderById
	err := db.Get(&tender, query, tenderId)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса на получение тендера: %v", err)
	}
	return &tender, nil
}
func GetTenderByIdFromTenderVersions(db *sqlx.DB, tenderId string, version string) (*models.TenderById, error) {
	var tenderVersion models.TenderById
	err := db.Get(&tenderVersion, `
		SELECT name, description, service_type, status, version, created_at, updated_at
		FROM tender_versions
		WHERE tender_id = $1 AND version = $2
	`, tenderId, version)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения данных тендера: %v", err)
	}
	return &tenderVersion, nil
}

// Функция получения тендера с минимальным набором данных
func GetTender(db *sqlx.DB, tenderId string, username string) (*models.Tender, int, error) {
	tender, validStatus, validErr := ValidTenderPermission(db, username, tenderId)
	if validErr != nil {
		return nil, validStatus, fmt.Errorf("ошибка проверки доступа к тендеру: %v", validErr)
	}

	// 4. Возвращаем тендер, скрывая некоторые данные (если нужно)
	return &models.Tender{
		Id:          tender.Id,
		Name:        tender.Name,
		Description: tender.Description,
		Status:      tender.Status,
		ServiceType: tender.ServiceType,
		Version:     tender.Version,
		CreatedAt:   tender.CreatedAt,
	}, http.StatusOK, nil
}

func UpdateTenderStatus(db *sqlx.DB, tenderId string, status string) (*models.Tender, error) {
	query := `
		UPDATE tender
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING id, name, description, status, service_type, version, created_at, updated_at
	`
	var updatedTender models.Tender
	err := db.Get(&updatedTender, query, status, tenderId)
	if err != nil {
		return nil, fmt.Errorf("ошибка при обновлении статуса тендера: %v", err)
	}

	return &updatedTender, nil
}

func SaveTenderVersion(db *sqlx.DB, tenderId string, currentTender *models.TenderById) (int, error)  {
	_, err:= db.Exec(`
		INSERT INTO tender_versions (tender_id, name, description, service_type, status, version, organization_id, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, tenderId, currentTender.Name, currentTender.Description, currentTender.ServiceType, currentTender.Status, currentTender.Version, currentTender.OrganizationId, currentTender.UserId, currentTender.CreatedAt, currentTender.UpdatedAt)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("ошибка сохранения текущей версии тендера: %v", err)
	}
	return http.StatusOK, nil
}

func UpdateTender(db *sqlx.DB, tenderId string, username string, newTenderData *models.Tender) (*models.Tender, int, error) {
	fmt.Println("UpdateTender username", username)

	currentTender, validStatus, validErr := ValidTenderPermission(db, username, tenderId)
	if validErr != nil {
		return nil, validStatus, fmt.Errorf("ошибка проверки доступа к тендеру: %v", validErr)
	}
	saveStatus, saveErr := SaveTenderVersion(db, tenderId, currentTender)
	if saveErr != nil {
		return nil, saveStatus, saveErr
	}

	query := `UPDATE tender SET `
	args := []interface{}{}
	counter := 1 // Счетчик для порядковых номеров параметров ($1, $2 и т.д.)

	// Если передан параметр name
	if newTenderData.Name != "" {
		query += fmt.Sprintf("name = $%d, ", counter)
		args = append(args, newTenderData.Name)
		counter++
	}

	// Если передан параметр description
	if newTenderData.Description != "" {
		query += fmt.Sprintf("description = $%d, ", counter)
		args = append(args, newTenderData.Description)
		counter++
	}

	// Если передан параметр service_type
	if newTenderData.ServiceType != "" {
		query += fmt.Sprintf("service_type = $%d, ", counter)
		args = append(args, newTenderData.ServiceType)
		counter++
	}
	newVersion := currentTender.Version + 1
	// Добавляем инкремент версии и время обновления
	query += fmt.Sprintf("version = %d, updated_at = CURRENT_TIMESTAMP WHERE id = $%d", newVersion, counter)
	args = append(args, tenderId)

	// Выполнение запроса
	_, err := db.Exec(query, args...)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("ошибка обновления тендера: %v", err)
	}

	// Получение обновленного тендера
	updatedTender, tenderStatus, err := GetTender(db, tenderId, username)
	if err != nil {
		return nil, tenderStatus, fmt.Errorf("ошибка получения обновленного тендера: %v", err)
	}

	return updatedTender, http.StatusOK, nil

}

func RollbackTenderVersion(db *sqlx.DB, tenderId string, version string, username string) (*models.Tender, int, error) {
	tender, validStatus, validErr := ValidTenderPermission(db, username, tenderId)
	if validErr != nil {
		return nil, validStatus, fmt.Errorf("ошибка проверки доступа к тендеру: %v", validErr)
	}
	// Получить нужную версию из таблицы tender_versions
	tenderVersion, err := GetTenderByIdFromTenderVersions(db, tenderId, version)
	if err != nil {
		return nil, http.StatusBadRequest,fmt.Errorf("ошибка получения тендера: %v", err)
	}
	// Сохранить текущий тендер в tender_versions
	saveStatus, saveErr := SaveTenderVersion(db, tenderId, tender)
	if saveErr != nil {
		return nil, saveStatus, saveErr
	}
	// Обновить текущий тендер, установив данные из старой версии
	_, err = db.Exec(`
		UPDATE tender
		SET name = $1, description = $2, service_type = $3, status = $4, version = version + 1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`, tenderVersion.Name, tenderVersion.Description, tenderVersion.ServiceType, tenderVersion.Status, tenderId)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("ошибка отката тендера к версии %s: %v", version, err)
	}

	// Вернуть обновленный тендер
	updatedTender, userStatus, err := GetTender(db, tenderId, username)
	if err != nil {
		return nil, userStatus, fmt.Errorf("ошибка получения откатанного тендера: %v", err)
	}

	return updatedTender, http.StatusOK, nil
}


