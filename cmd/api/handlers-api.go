package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"stripe-app/internal/cards"

	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type stripePayload struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Content string `json:"content,omitempty"`
	ID      int    `json:"id,omitempty"`
}

func (app *application) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	var payload stripePayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	amount, err := strconv.Atoi(payload.Amount)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	card := cards.Card{
		Secret:   app.config.stripe.secret,
		Key:      app.config.stripe.key,
		Currency: payload.Currency,
	}

	okay := true

	pi, msg, err := card.Charge(payload.Currency, amount)
	if err != nil {
		okay = false
	}

	if okay {
		out, err := json.MarshalIndent(pi, "", "   ")
		if err != nil {
			app.errorLog.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	} else {
		j := jsonResponse{
			OK:      false,
			Message: msg,
			Content: "",
		}

		out, err := json.MarshalIndent(j, "", "   ")
		if err != nil {
			app.errorLog.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	}
}
func (app *application) GetWidgetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fmt.Println("Raw ID from URL:", id) // Debugging output

	widgetID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Invalid ID (not a number):", id) // Debugging output
		http.Error(w, `{"error": "Invalid widget ID"}`, http.StatusBadRequest)
		return
	}

	fmt.Println("Fetching widget with ID:", widgetID) // Debugging output
	widget, err := app.DB.GetWidget(widgetID)
	if err == sql.ErrNoRows {
		fmt.Println("No widget found with ID:", widgetID)
		http.Error(w, `{"error": "Widget not found"}`, http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("Database error:", err)
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	out, err := json.MarshalIndent(widget, "", "   ")
	if err != nil {
		fmt.Println("JSON encoding error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
