package config

type Payment struct {
	PaymentMethod string  `json:"paymentMethod"`
	UserID        int     `json:"userId"`
	OrderID       int     `json:"orderId"`
	Total         float64 `json:"total"`
}

type PaymentPageData struct {
	Client *JWTtoken
	Order  *Order
	Total  float64
}
