package config

import "time"

type Order struct {
	ID          int         `json:"id"`
	UserID      int         `json:"userId"`
	TableNumber int         `json:"tableNumber"`
	Complete    bool        `json:"complete"`
	CreatedAt   time.Time   `json:"createdAt"`
	Items       []OrderItem `json:"items"`
}

type OrderItem struct {
	OrderID     int    `json:"orderId,omitempty"`
	ItemID      int    `json:"itemId"`
	Name        string `json:"name"`
	Quantity    int    `json:"quantity"`
	Instruction string `json:"instruction"`
}
