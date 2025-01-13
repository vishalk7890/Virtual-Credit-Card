package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

// Models is the wrapper for all models
type Models struct {
	DB DBModel
}

// NewModels returns a model type with database connection pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

// Widget is the type for all widgets
type Widget struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	InventoryLevel int       `json:"inventory_level"`
	Price          int       `json:"price"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

func (m *DBModel) GetWidget(id int) (Widget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var widget Widget

	query := m.DB.QueryRowContext(ctx, "SELECT * FROM widgets WHERE id = ?", id)
	err := query.Scan(&widget.ID, &widget.Name, &widget.Description, &widget.InventoryLevel, &widget.Price, &widget.CreatedAt, &widget.UpdatedAt)
	if err != nil {
		return widget, fmt.Errorf("error querying one widget: %w", err)
	}
	return widget, nil
}
