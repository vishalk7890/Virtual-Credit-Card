package main

import (
	"fmt"
	"net/http"
	"strconv"
	"stripe-app/internal/cards"
	"stripe-app/internal/models"
	"time"

	"github.com/go-chi/chi/v5"
)

// VirtualTerminal displays the virtual terminal page
func (app *application) VirtualTerminal(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "terminal", &templateData{}, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "home", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

// PaymentSucceeded displays the receipt page
func (app *application) PaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// read posted data
	//cardHolder := r.Form.Get("cardholder_name")
	email := r.Form.Get("email")
	first_name := r.Form.Get("first_name")
	last_name := r.Form.Get("last_name")
	paymentIntent := r.Form.Get("payment_intent")
	paymentMethod := r.Form.Get("payment_method")
	paymentAmount := r.Form.Get("payment_amount")
	paymentCurrency := r.Form.Get("payment_currency")
	widgetId, _ := strconv.Atoi(r.Form.Get("product_id"))

	card := cards.Card{
		Secret: app.config.stripe.secret,
		Key:    app.config.stripe.key,
	}

	// get payment method
	pi, err := card.RetrieveGetPaymentIntent(paymentIntent)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	pm, err := card.GetPaymentMethod(paymentMethod)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	lastFour := pm.Card.Last4
	expiryMonth := strconv.Itoa(int(pm.Card.ExpMonth))
	expiryYear := strconv.Itoa(int(pm.Card.ExpYear))

	// craete a new customer
	customerID, err := app.SaveCustomer(first_name, last_name, email)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	app.infoLog.Println("Customer ID:", customerID)

	// createa new transaction
	amount, err := strconv.Atoi(paymentAmount)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	txn := models.Transaction{
		Amount:              amount,
		Currency:            paymentCurrency,
		LastFour:            lastFour,
		ExpiryMonth:         expiryMonth,
		ExpiryYear:          expiryYear,
		BankReturnCode:      pi.Charges.Data[0].ID,
		TransactionStatusID: 2,
	}
	txnId, err := app.SaveTransaction(txn)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.infoLog.Println("Transaction ID:", txnId)

	//craete a new order

	order := models.Order{
		WidgetID:      widgetId,
		TransactionID: txnId,
		CustomerID:    customerID,
		StatusID:      1,
		Quantity:      1,
		Amount:        amount,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err = app.SaveOrder(order)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	data := make(map[string]interface{})
	//data["cardholder"] = cardHolder
	data["email"] = email
	data["pi"] = paymentIntent
	data["pm"] = paymentMethod
	data["pa"] = paymentAmount
	data["pc"] = paymentCurrency
	data["last_four"] = lastFour
	data["expiry_month"] = expiryMonth
	data["expiry_year"] = expiryYear
	data["bank_return_code"] = pi.Charges.Data[0].ID
	data["first_name"] = first_name
	data["last_name"] = last_name


	// should write this data to session, and then redirect user to new page

	if err := app.renderTemplate(w, r, "succeeded", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}
}

// ChargeOnce displays the page to buy one widget
func (app *application) ChargeOnce(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fmt.Println("Raw ID from URL:", id) // Debugging output

	widgetID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Invalid ID (not a number):", id) // Debugging output
		http.Error(w, `{"error": "Invalid widget ID"}`, http.StatusBadRequest)
		return
	}
	widget, err := app.DB.GetWidget(widgetID)
	if err != nil {
		fmt.Println("Error fetching widget:", err) // Debugging output
		return
	}

	data := make(map[string]interface{})
	data["widget"] = widget

	if err := app.renderTemplate(w, r, "buy-once", &templateData{
		Data: data,
	}, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) SaveCustomer(firstName, lastName, email string) (int, error) {
	customer := models.Customer{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}
	id, err := app.DB.InsertCustomer(customer)
	if err != nil {
		return 0, err

	}
	return id, nil
}

func (app *application) SaveTransaction(txn models.Transaction) (int, error) {
	id, err := app.DB.InsertTransaction(txn)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (app *application) SaveOrder(order models.Order) (int, error) {
	id, err := app.DB.InsertOrder(order)
	if err != nil {
		return 0, err
	}
	return id, nil
}
