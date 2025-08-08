package model

import (
	"database/sql"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type ModelConnection struct {
	DB *sql.DB
}

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JWTtoken struct {
	ID      int
	Name    string
	IsAdmin bool
	IsCheff bool
	jwt.RegisteredClaims
}

type Order struct {
	ID          int         `json:"id"`
	UserID      int         `json:"userId"`
	TableNumber int         `json:"tableNumber"`
	Complete    bool        `json:"complete"`
	CreatedAt   time.Time   `json:"createdAt"`
	Items       []OrderItem `json:"items"`
}

// OrderItem represents a single item within an order.
type OrderItem struct {
	OrderID     int    `json:"orderId,omitempty"`
	ItemID      int    `json:"itemId"`
	Quantity    int    `json:"quantity"`
	Instruction string `json:"instruction"`
}

type Payment struct {
	PaymentMethod string `json:"paymentMethod"`
	UserID        int    `json:"userId"`
	OrderID       int    `json:"orderId"`
	Total         int    `json:"total"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Item struct {
	ID          int     `json:"id"`
	CategoryID  int     `json:"categoryId"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"` // Use float64 for price
	Description string  `json:"description"`
}
