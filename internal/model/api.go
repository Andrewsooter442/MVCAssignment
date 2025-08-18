package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Andrewsooter442/MVCAssignment/types"
)

func (model *ModelConnection) PlaceOrder(order *types.Order) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := model.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

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

func (model *ModelConnection) CompletePayment(paymentDetails *types.Payment) error {
	query := `
        INSERT INTO payment (order_id, user_id, total, tip, paid, method)
        VALUES (?, ?, ?, 0, TRUE, ? )
    `
	_, err := model.DB.Exec(query, paymentDetails.OrderID, paymentDetails.UserID, paymentDetails.Total, paymentDetails.PaymentMethod)
	if err != nil {
		log.Printf("Error executing insert statement for new user: %v", err)
		return errors.New("failed to create user in database")
	}

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

func (model *ModelConnection) CompleteOrderItem(orderID int, itemID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE order_items SET complete = TRUE WHERE order_id = ? AND item_id = ?`
	res, err := model.DB.ExecContext(ctx, query, orderID, itemID)
	if err != nil {
		log.Printf("Error updating order item completion status: %v", err)
		return fmt.Errorf("could not update item status")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify update")
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no matching item found in order to complete")
	}

	return nil
}

func (model *ModelConnection) CheckOrderItemByOrderID(id int) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM order_items WHERE order_id = ? AND complete = ? LIMIT 1)`

	err := model.DB.QueryRow(query, id, false).Scan(&exists)
	if err != nil {
		fmt.Println("Error checking order item:", err)
		return false
	}

	return exists
}

func (model *ModelConnection) GetUserIDForOrder(orderID int) (int, error) {
	var userID int
	query := `SELECT user_id FROM orders WHERE id = ?`
	err := model.DB.QueryRow(query, orderID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no order found with ID %d", orderID)
		}
		return 0, err
	}
	return userID, nil
}

func (model *ModelConnection) GetLatestUnpaidOrderForUser(userID int) (*types.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var latestOrderID int
	query := `
        SELECT o.id
        FROM orders o
        LEFT JOIN payment p ON o.id = p.order_id
        WHERE o.user_id = ? AND p.order_id IS NULL
        ORDER BY o.created_at DESC
        LIMIT 1;
    `

	err := model.DB.QueryRowContext(ctx, query, userID).Scan(&latestOrderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Error finding latest unpaid order ID for user %d: %v", userID, err)
		return nil, fmt.Errorf("could not retrieve latest order ID: %w", err)
	}

	return model.GetOrderById(latestOrderID)
}
