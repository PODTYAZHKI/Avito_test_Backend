package models

type Bid struct {
	Id          string `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"-"`
	Status      string `db:"status" json:"status"`
	TenderId    string `db:"tender_id" json:"-"`
	AuthorType  string `db:"author_type" json:"authorType"`
	AuthorId    string `db:"author_id" json:"authorId"`
	Version     int    `db:"version" json:"version"`
	CreatedAt   string `db:"created_at" json:"createdAt"`
	UpdatedAt   string `db:"updated_at" json:"-"`
}

type BidById struct {
	Id          string `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Status      string `db:"status"`
	TenderId    string `db:"tender_id"`
	AuthorType  string `db:"author_type"`
	AuthorId    string `db:"author_id"`
	Version     int    `db:"version"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
}

type BidRequest struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	Status          string `json:"status"`
	TenderId        string `json:"tenderId"`
	OrganizationId  string `json:"organizationId"`
	CreatorUsername string `json:"creatorUsername"`
}
