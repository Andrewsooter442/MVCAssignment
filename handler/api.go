package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/Andrewsooter442/MVCAssignment/types"
)

func (app *Application) HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   false,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Application) HandlePlaceOrder(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(types.UserObject).(*types.JWTtoken)
	if !ok {
		http.Error(w, "Could not retrieve user claims", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	if err := validateOrderForm(r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var order types.Order
	order.UserID = claims.ID
	order.TableNumber, _ = strconv.Atoi(r.FormValue("tableNumber"))

	itemIDs := r.Form["itemId"]
	quantities := r.Form["quantity"]
	instructions := r.Form["instruction"]

	for i := 0; i < len(itemIDs); i++ {
		itemID, _ := strconv.Atoi(itemIDs[i])
		quantity, _ := strconv.Atoi(quantities[i])
		order.Items = append(order.Items, types.OrderItem{
			ItemID:      itemID,
			Quantity:    quantity,
			Instruction: instructions[i],
		})
	}

	_, err := app.Pool.PlaceOrder(&order)
	if err != nil {
		http.Error(w, "Failed to place order", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/api/payment", http.StatusSeeOther)
}

func (app *Application) HandleCompleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid Order ID in URL", http.StatusBadRequest)
		return
	}

	err = app.Pool.CompleteOrder(orderID)
	if err != nil {
		http.Error(w, "Failed to complete order", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func validateOrderForm(form url.Values) error {
	itemIDs := form["itemId"]
	quantities := form["quantity"]
	instructions := form["instruction"]

	if len(itemIDs) == 0 {
		return errors.New("order must contain at least one item")
	}

	if len(itemIDs) != len(quantities) || len(itemIDs) != len(instructions) {
		return errors.New("mismatched number of order items, quantities, and instructions")
	}

	tableNumStr := form.Get("tableNumber")
	if tableNumStr == "" {
		return errors.New("tableNumber is a required field")
	}

	tableNumber, err := strconv.Atoi(tableNumStr)
	if err != nil || tableNumber <= 0 {
		return errors.New("invalid table number")
	}

	return nil
}

func validatePaymentData(form url.Values) error {
	requiredFields := []string{"paymentMethod", "orderId", "total"}
	for _, field := range requiredFields {
		if form.Get(field) == "" {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	if _, err := strconv.Atoi(form.Get("orderId")); err != nil {
		return errors.New("invalid format for orderId, must be an integer")
	}

	total, err := strconv.ParseFloat(form.Get("total"), 64)
	if err != nil {
		return errors.New("invalid format for total, must be a number")
	}
	if total <= 0 {
		return errors.New("invalid value for total, must be a positive number")
	}

	method := form.Get("paymentMethod")
	allowedMethods := map[string]bool{
		"Monero":      true,
		"Credit Card": true,
		"Debit Card":  true,
		"Bitcoin":     true,
		"Cash":        true,
	}

	if _, ok := allowedMethods[method]; !ok {
		return fmt.Errorf("payment method '%s' is not supported", method)
	}
	return nil
}

func (app *Application) HandleGetPayment(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(types.UserObject).(*types.JWTtoken)
	if !ok {
		http.Error(w, "Could not retrieve user claims", http.StatusInternalServerError)
		return
	}

	order, err := app.Pool.GetLatestUnpaidOrderForUser(claims.ID)
	if err != nil {
		log.Printf("Database error in HandleGetPayment: %v", err)
		http.Error(w, "Could not retrieve order details", http.StatusInternalServerError)
		return
	}

	var total float64
	if order != nil {
		for _, item := range order.Items {
			total += item.Price * float64(item.Quantity)
		}
	}

	pageData := types.PaymentPageData{
		Client: claims,
		Order:  order,
		Total:  total,
	}

	err = types.Templates.ExecuteTemplate(w, "payment.html", pageData)
	if err != nil {
		log.Printf("Template execution error in HandleGetPayment: %v", err)
		http.Error(w, "Failed to render the payment page.", http.StatusInternalServerError)
	}
}

func (app *Application) HandlePostPayment(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(types.UserObject).(*types.JWTtoken)
	if !ok {
		http.Error(w, "Could not retrieve user claims", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	if err := validatePaymentData(r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderID, _ := strconv.Atoi(r.FormValue("orderId"))
	total, _ := strconv.ParseFloat(r.FormValue("total"), 64)
	paymentMethod := r.FormValue("paymentMethod")

	orderUserID, err := app.Pool.GetUserIDForOrder(orderID)
	if err != nil {
		http.Error(w, "Invalid Order ID", http.StatusNotFound)
		return
	}
	if orderUserID != claims.ID {
		http.Error(w, "Forbidden: You cannot pay for an order that is not yours.", http.StatusForbidden)
		return
	}

	paymentDetails := &types.Payment{
		OrderID:       orderID,
		UserID:        claims.ID,
		Total:         total,
		PaymentMethod: paymentMethod,
	}

	if err := app.Pool.CompletePayment(paymentDetails); err != nil {
		log.Printf("Failed to complete payment: %v", err)
		http.Error(w, "Failed to process payment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) HandleCompleteOrderItem(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		OrderID int `json:"orderId"`
		ItemID  int `json:"itemId"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payload.OrderID == 0 || payload.ItemID == 0 {
		http.Error(w, "Missing orderId or itemId", http.StatusBadRequest)
		return
	}

	err = app.Pool.CompleteOrderItem(payload.OrderID, payload.ItemID)
	if err != nil {
		http.Error(w, "Failed to update item status", http.StatusInternalServerError)
		return
	}

	isEmpty := app.Pool.CheckOrderItemByOrderID(payload.OrderID)
	fmt.Println(isEmpty)

	if !isEmpty {
		app.Pool.CompleteOrder(payload.OrderID)

	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
