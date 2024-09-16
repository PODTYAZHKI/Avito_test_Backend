package store

import (
	"fmt"
	"net/http"
	"tender-service/models"

	"github.com/jmoiron/sqlx"
)

func GetBid(db *sqlx.DB, bidId string) (*models.Bid, error) {	
	query := `
        SELECT name, description, status, version, tender_id, author_type, author_id, created_at, updated_at
        FROM bid
        WHERE id = $1
    `
	var bid models.Bid
	err := db.Get(&bid, query, bidId)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса на получение предложения: %v", err)
	}
	return &bid, nil
}


func GetBidById(db *sqlx.DB, bidId string) (*models.BidById, error) {	
	query := `
        SELECT name, description, status, version, tender_id, author_type, author_id, created_at, updated_at
        FROM bid
        WHERE id = $1
    `
	var bid models.BidById
	err := db.Get(&bid, query, bidId)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса на получение предложения: %v", err)
	}
	return &bid, nil
}

func CreateBid(db *sqlx.DB, bid *models.Bid) (*models.Bid, error) {
	query := `
		INSERT INTO bid (name, description, status, version, tender_id, author_type, author_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := db.QueryRowx(query, bid.Name, bid.Description, bid.Status, bid.Version, bid.TenderId, bid.AuthorType, bid.AuthorId).StructScan(bid)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании предложения: %v", err)
	}

	return bid, nil
}

func GetUserBids(db *sqlx.DB, organizationId string, userId string, limit int, offset int) (*[]models.Bid, error) {
	fmt.Println("organizationId", organizationId)
	fmt.Println("userId", userId)
	query := `
        SELECT id, name, description, status, version, tender_id, author_type, author_id, created_at
        FROM bid
        WHERE (author_type = 'User' AND author_id = $1)
        OR (author_type = 'Organization' AND author_id = $2)
        ORDER BY name
        LIMIT $3 OFFSET $4
    `
	var bids []models.Bid

	err := db.Select(&bids, query, userId, organizationId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса на получение предложений: %v", err)
	}

	return &bids, nil
}

func GetTenderBids(db *sqlx.DB, tenderId string, limit int, offset int) (*[]models.Bid, error) {
	query := `
        SELECT id, name, description, status, version, tender_id, author_type, author_id, created_at
        FROM bid
        WHERE tender_id = $1
        ORDER BY name
        LIMIT $2 OFFSET $3
    `
	var bids []models.Bid

	err := db.Select(&bids, query, tenderId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса на получение предложений тендера: %v", err)
	}

	return &bids, nil
}

func GetBidStatus(db *sqlx.DB, bidId string) (string, error) {
	query := `
        SELECT status
        FROM bid
        WHERE id = $1
    `
	var status string
	err := db.Get(&status, query, bidId)
	if err != nil {
		return status, fmt.Errorf("ошибка получения статуса предложения: %v", err)
	}

	return status, nil
}

func UpdateBidStatus(db *sqlx.DB, bidId string, status string) (*models.Bid, error) {
	query := `
		UPDATE bid
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING id, name, status, version, author_type, author_id, created_at
	`
	var updatedBid models.Bid
	err := db.Get(&updatedBid, query, status, bidId)
	if err != nil {
		return nil, fmt.Errorf("ошибка при обновлении статуса тендера: %v", err)
	}

	return &updatedBid, nil
}

func SaveBidVersion(db *sqlx.DB, bidId string, currentBid *models.BidById) (int, error)  {
	fmt.Println("currentBid", currentBid)
	fmt.Println("bidId", bidId)
	_, err:= db.Exec(`
		INSERT INTO bid_versions (bid_id ,name, description, status, tender_id, author_id,
		author_type,
		version,
		created_at,
        updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, bidId, currentBid.Name, currentBid.Description, currentBid.Status, currentBid.TenderId, currentBid.AuthorId, currentBid.AuthorType, currentBid.Version, currentBid.CreatedAt, currentBid.UpdatedAt)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("ошибка сохранения текущей версии предложения: %v", err)
	}
	return http.StatusOK, nil
}

func UpdateBid(db *sqlx.DB, bidId string, username string, newBidData *models.Bid) (*models.Bid, int, error) {

	currentBid, currentBidErr := GetBidById(db, bidId)
	if currentBidErr != nil {
		return nil, http.StatusBadRequest, currentBidErr
	}
	saveStatus, saveErr := SaveBidVersion(db, bidId, currentBid)
	if saveErr != nil {
		return nil, saveStatus, saveErr
	}

	query := `UPDATE bid SET `
	args := []interface{}{}
	counter := 1 // Счетчик для порядковых номеров параметров ($1, $2 и т.д.)

	// Если передан параметр name
	if newBidData.Name != "" {
		query += fmt.Sprintf("name = $%d, ", counter)
		args = append(args, newBidData.Name)
		counter++
	}

	// Если передан параметр description
	if newBidData.Description != "" {
		query += fmt.Sprintf("description = $%d, ", counter)
		args = append(args, newBidData.Description)
		counter++
	}

	// Если передан параметр service_type
	newVersion := currentBid.Version + 1
	// Добавляем инкремент версии и время обновления
	query += fmt.Sprintf("version = %d, updated_at = CURRENT_TIMESTAMP WHERE id = $%d", newVersion, counter)
	args = append(args, bidId)
	fmt.Println("query", query)
	// Выполнение запроса
	_, err := db.Exec(query, args...)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("ошибка обновления предложения: %v", err)
	}

	// Получение обновленного тендера
	updatedBid, bidErr := GetBid(db, bidId)
	if bidErr != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("ошибка получения обновленного предложения: %v", err)
	}

	return updatedBid, http.StatusOK, nil

}

func GetBidByIdFromBidVersions(db *sqlx.DB, bidId string, version string) (*models.BidById, error) {
	var bidVersion models.BidById
	err := db.Get(&bidVersion, `
		SELECT name, description, status, version, tender_id, author_type, author_id, created_at, updated_at
        FROM bid_versions
        WHERE bid_id = $1 AND version = $2
	`, bidId, version)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения данных предложения: %v", err)
	}
	return &bidVersion, nil
}

func RollbackBidVersion(db *sqlx.DB, bidId string, version string, username string) (*models.Bid, int, error) {

	currentBid, currentBidErr := GetBidById(db, bidId)
	if currentBidErr != nil {
		return nil, http.StatusBadRequest, currentBidErr
	}

	bidVersion, err := GetBidByIdFromBidVersions(db, bidId, version)
	if err != nil {
		return nil, http.StatusBadRequest,fmt.Errorf("ошибка получения предложения: %v", err)
	}
	// Сохранить текущий тендер в tender_versions
	saveStatus, saveErr := SaveBidVersion(db, bidId, currentBid)
	if saveErr != nil {
		return nil, saveStatus, saveErr
	}
	// Обновить текущий тендер, установив данные из старой версии
	_, err = db.Exec(`
		UPDATE bid
		SET name = $1, description = $2, status = $3, version = version + 1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`, bidVersion.Name, bidVersion.Description, bidVersion.Status, bidId)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("ошибка отката предложения к версии %s: %v", version, err)
	}

	// Вернуть обновленный тендер
	updatedBid, err := GetBid(db, bidId)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("ошибка получения откатанного предложения: %v", err)
	}

	return updatedBid, http.StatusOK, nil
}