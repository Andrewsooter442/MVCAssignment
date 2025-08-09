package config

type Payment struct {
	PaymentMethod string `json:"paymentMethod"`
	UserID        int    `json:"userId"`
	OrderID       int    `json:"orderId"`
	Total         int    `json:"total"`
}
