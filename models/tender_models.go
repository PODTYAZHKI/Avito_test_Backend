package models

type Tender struct {
	Id             string `db:"id" json:"id"`
	Name           string `db:"name" json:"name"`
	Description    string `db:"description" json:"description"`
	Status         string `db:"status" json:"status"`
	OrganizationId string `db:"organization_id" json:"-"`
	UserId         string `db:"user_id" json:"-"`
	ServiceType    string `db:"service_type" json:"service_type"`
	Version        int    `db:"version" json:"version"`
	CreatedAt      string `db:"created_at" json:"createdAt"`
	UpdatedAt      string `db:"updated_at" json:"-"`
}

type TenderById struct {
	Id             string `db:"id"`
	Name           string `db:"name"`
	Description    string `db:"description"`
	Status         string `db:"status"`
	OrganizationId string `db:"organization_id"`
	UserId         string `db:"user_id"`
	ServiceType    string `db:"service_type"`
	Version        int    `db:"version"`
	CreatedAt      string `db:"created_at"`
	UpdatedAt      string `db:"updated_at"`
}

type TenderRequest struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	ServiceType     string `json:"serviceType"`
	Status          string `json:"status"`
	OrganizationId  string `json:"organizationId"`
	CreatorUsername string `json:"creatorUsername"`
}
