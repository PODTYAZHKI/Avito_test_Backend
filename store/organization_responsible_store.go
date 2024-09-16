package store

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func GetIdByOrganizationId(db *sqlx.DB, organizationId string) (string, error) {
	var organizationResponsibleId string
	err := db.Get(&organizationResponsibleId, "SELECT id FROM organization_responsible WHERE organization_id = $1", organizationId)
	if err != nil {
		return organizationResponsibleId, fmt.Errorf("ошибка получения id из таблицы organisation_responsible")
	}
	return organizationResponsibleId, nil
}

func GetIdByUserId(db *sqlx.DB, userId string) (string, error) {
	var organizationResponsibleId string
	err := db.Get(&userId, "SELECT id FROM organization_responsible WHERE user_id = $1", userId)
	if err != nil {
		return organizationResponsibleId, fmt.Errorf("ошибка получения id из таблицы organisation_responsible")
	}
	return organizationResponsibleId, nil
}

func GetOrganizationIdByUserId(db *sqlx.DB, userId string) (string, error) {
	var organizationId string
	err := db.Get(&organizationId, "SELECT organization_id FROM organization_responsible WHERE user_id = $1", userId)
	if err != nil {
		return organizationId, fmt.Errorf("ошибка получения id организации из таблицы organisation_responsible")
	}
	return organizationId, nil
}