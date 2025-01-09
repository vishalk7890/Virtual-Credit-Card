package models

import (
	"database/sql"
	"time"
)

type DBmodel struct {
	DB *sql.DB
}

type Models struct {
	DB DBmodel
}

func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBmodel{DB: db},
	}
}

type Widget struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	InventoryLevel int       `json:"inventorylevel"`
	Price          int       `json:"price"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}
