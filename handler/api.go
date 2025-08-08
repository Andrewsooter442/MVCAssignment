package handler

import (
	"errors"
	"fmt"
	"github.com/Andrewsooter442/MVCAssignment/internal/model"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const userObject = "user"

func (app *Application) HandleApiRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/api/"):]
	switch path {
	case "logout":
		app.HandleLogout(w, r)
	case "placeOrder":
		app.HandlePlaceOrder(w, r)
	case "completeOrderItem":

		app.HandleCompleteOrderItem(w, r)
	case "paymentDone":

	default:
		http.NotFound(w, r)
	}
}

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

func (app *Application) HandleCompleteOrderItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ItemID, err := strconv.Atoi(r.FormValue("itemId"))
	if err != nil {
		http.Error(w, "Invalid itemId", http.StatusBadRequest)
		return
	}

	err = app.Pool.CompleteOrder(ItemID)

}

func (app *Application) HandlePaymentDone(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(userObject).(*model.JWTtoken)
	if !ok {
		http.Error(w, "Could not retrieve user claims", http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	if err := validatePaymentData(r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var paymentDetails model.Payment
	var err error

	paymentDetails.UserID = claims.ID
	paymentDetails.OrderID, _ = strconv.Atoi(r.FormValue("orderId"))
	paymentDetails.PaymentMethod = r.FormValue("paymentMethod")
	paymentDetails.Total, _ = strconv.Atoi(r.FormValue("total"))

	err = app.Pool.CompletePayment(&paymentDetails)
	if err != nil {
		http.Error(w, "Failed to place order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Order placed successfully! Order ID: %d", paymentDetails.OrderID)

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

	_, err := strconv.Atoi(form.Get("total"))
	if err != nil {
		return errors.New("invalid format for total, must be an integer")
	}

	return nil
}

func (app *Application) HandlePlaceOrder(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(userObject).(*model.JWTtoken)
	if !ok {
		http.Error(w, "Could not retrieve user claims", http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

	var order model.Order
	var err error

	order.UserID = claims.ID

	order.TableNumber, _ = strconv.Atoi(r.FormValue("tableNumber"))

	itemIDs := r.Form["itemId"]
	quantities := r.Form["quantity"]
	instructions := r.Form["instruction"]

	for i := 0; i < len(itemIDs); i++ {
		itemID, _ := strconv.Atoi(itemIDs[i])
		quantity, _ := strconv.Atoi(quantities[i])

		order.Items = append(order.Items, model.OrderItem{
			ItemID:      itemID,
			Quantity:    quantity,
			Instruction: instructions[i],
		})
	}

	fmt.Println(order, "coming from handler/api.go")

	orderID, err := app.Pool.PlaceOrder(&order)
	if err != nil {
		http.Error(w, "Failed to place order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Order placed successfully! Order ID: %d", orderID)
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
