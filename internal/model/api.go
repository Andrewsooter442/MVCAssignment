package model

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

func (model *ModelConnection) PlaceOrder(order *Order) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := model.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Step 1: Insert the main order record.
	orderQuery := `
        INSERT INTO orders (user_id, table_no, complete, created_at)
        VALUES (?, ?, ?, UTC_TIMESTAMP())
    `
	res, err := tx.ExecContext(ctx, orderQuery, order.UserID, order.TableNumber, order.Complete)
	if err != nil {
		log.Printf("Error executing insert statement for new order: %v", err)
		return 0, fmt.Errorf("failed to insert order: %w", err)
	}

	orderID, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID for order: %v", err)
		return 0, fmt.Errorf("failed to get order ID: %w", err)
	}

	itemStmt, err := tx.PrepareContext(ctx, `
        INSERT INTO order_items (order_id, item_id, quantity, instruction)
        VALUES (?, ?, ?, ?)
    `)
	if err != nil {
		log.Printf("Error preparing statement for order items: %v", err)
		return 0, fmt.Errorf("failed to prepare item statement: %w", err)
	}
	defer itemStmt.Close()

	for _, item := range order.Items {
		_, err := itemStmt.ExecContext(ctx, orderID, item.ItemID, item.Quantity, item.Instruction)
		if err != nil {
			log.Printf("Error executing insert for order item %d: %v", item.ItemID, err)
			return 0, fmt.Errorf("failed to insert order item %d: %w", item.ItemID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("could not commit transaction: %w", err)
	}

	return orderID, nil
}

func (model *ModelConnection) CompletePayment(paymentDetails *Payment) error {
	query := `
        INSERT INTO payment (order_id, user_id, total, tip, paid, method)
        VALUES (?, ?, ?, 0, TRUE, ? )
    `
	_, err := model.DB.Exec(query, paymentDetails.OrderID, paymentDetails.UserID, paymentDetails.Total, paymentDetails.PaymentMethod)
	if err != nil {
		log.Printf("Error executing insert statement for new user: %v", err)
		return errors.New("failed to create user in database")
	}

	//fmt.Println("User created successfully")
	return nil

}

func (model *ModelConnection) CompleteOrder(orderID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := model.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	var incompleteItemsCount int
	query := `SELECT COUNT(*) FROM order_items WHERE order_id = ? AND complete = FALSE`
	err = tx.QueryRowContext(ctx, query, orderID).Scan(&incompleteItemsCount)
	if err != nil {
		log.Printf("Error checking for incomplete items for order ")
		return fmt.Errorf("could not check order items: %w", err)
	}

	updateQuery := `UPDATE orders SET complete = TRUE WHERE id = ?`
	_, err = tx.ExecContext(ctx, updateQuery, orderID)
	if err != nil {
		log.Printf("Error executing update for completing order")
		return fmt.Errorf("failed to update order: %w", err)
	}

	return tx.Commit()
}
