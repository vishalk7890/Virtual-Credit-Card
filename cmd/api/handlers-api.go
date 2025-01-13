package main

import (
	"database/sql"
	"encoding/json"
	"stripe-app/internal/cards"

	"net/http"
	"strconv"
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
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		app.errorLog.Println("Missing id parameter")
		return
	}

	id, err := strconv.Atoi(idStr)
	app.infoLog.Printf("URL Parameter 'id': %q", id)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		app.errorLog.Println("Invalid id parameter:", err)
		return
	}

	widget, err := app.DB.GetWidget(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Widget not found", http.StatusNotFound)
			app.errorLog.Println("Widget not found:", err)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			app.errorLog.Println("Failed to get widget:", err)
		}
		return
	}

	out, err := json.MarshalIndent(widget, "", "   ")
	if err != nil {
		http.Error(w, "Failed to marshal widget", http.StatusInternalServerError)
		app.errorLog.Println("Failed to marshal widget:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
