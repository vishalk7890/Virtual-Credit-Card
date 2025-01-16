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
	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var widget Widget

	rows := m.DB.QueryRow("SELECT COUNT(*) FROM widgets")

	defer rows.Close()

	// Iterate over the rows and print the results
	for rows.Next() {
		var widget Widget
		err := rows.Scan(&widget.ID, &widget.Name)
		if err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			return widget, err
		}
		// Print each widget's data
		fmt.Printf("Widget ID: %d, Name: %s\n", widget.ID, widget.Name)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		fmt.Printf("Error iterating over rows: %v\n", err)

		return widget, nil
	}
	return widget, nil
}
