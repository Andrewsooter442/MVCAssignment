package model

import (
	"database/sql"
	"github.com/Andrewsooter442/MVCAssignment/types"
	"log"
)

func (model *ModelConnection) GetItemPrice(itemId int) (float64, error) {
	var price float64
	query := `SELECT price FROM items WHERE id = ?`
	err := model.DB.QueryRow(query, itemId).Scan(&price)
	if err != nil {
		log.Printf("Failed to get item price for item %d: %v", itemId, err)
		return 0, err
	}
	return price, nil
}

func (model *ModelConnection) GetOrderById(orderID int) (*types.Order, error) {

	var order types.Order

	orderQuery := `SELECT id, user_id, table_no, complete, created_at FROM orders WHERE id = ?`
	err := model.DB.QueryRow(orderQuery, orderID).Scan(&order.ID, &order.UserID, &order.TableNumber, &order.Complete, &order.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No order found with ID: %d", orderID)
			return nil, err
		}
		log.Printf("Error querying for order %d: %v", orderID, err)
		return nil, err
	}

	itemsQuery := `
		SELECT oi.item_id, i.name, oi.quantity, oi.instruction
		FROM order_items oi
		JOIN items i ON oi.item_id = i.id
		WHERE oi.order_id = ?
	`
	rows, err := model.DB.Query(itemsQuery, orderID)
	if err != nil {
		log.Printf("Error querying order items for order %d: %v", orderID, err)
		return nil, err
	}
	defer rows.Close()

	var orderItems []types.OrderItem
	for rows.Next() {
		var item types.OrderItem
		item.OrderID = orderID
		if err := rows.Scan(&item.ItemID, &item.Name, &item.Quantity, &item.Instruction); err != nil {
			log.Printf("Error scanning order item: %v", err)
			return nil, err
		}
		item.Price, err = model.GetItemPrice(item.ItemID)
		orderItems = append(orderItems, item)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Rows iteration error on order items: %v", err)
		return nil, err
	}

	order.Items = orderItems

	return &order, nil
}

func (model *ModelConnection) GetIncompleteOrders() ([]types.Order, error) {
	ordersQuery := `SELECT id, user_id, table_no, created_at FROM orders WHERE complete = FALSE ORDER BY created_at ASC`
	rows, err := model.DB.Query(ordersQuery)
	if err != nil {
		log.Printf("ChefView: Error querying incomplete orders: %v", err)
		return nil, err
	}
	defer rows.Close()

	orderMap := make(map[int]*types.Order)
	var orderedList []*types.Order

	for rows.Next() {
		var order types.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.TableNumber, &order.CreatedAt); err != nil {
			log.Printf("ChefView: Error scanning incomplete order: %v", err)
			return nil, err
		}
		order.Items = []types.OrderItem{}
		orderMap[order.ID] = &order
		orderedList = append(orderedList, &order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	if len(orderMap) == 0 {
		return []types.Order{}, nil
	}

	itemsQuery := `
		SELECT oi.order_id, oi.item_id, i.name, oi.quantity, oi.instruction
		FROM order_items oi
		JOIN items i ON oi.item_id = i.id
		WHERE oi.order_id IN (SELECT id FROM orders WHERE complete = FALSE) AND oi.complete = FALSE
	`
	itemRows, err := model.DB.Query(itemsQuery)
	if err != nil {
		log.Printf("ChefView: Error querying items for incomplete orders: %v", err)
		return nil, err
	}
	defer itemRows.Close()

	for itemRows.Next() {
		var item types.OrderItem
		if err := itemRows.Scan(&item.OrderID, &item.ItemID, &item.Name, &item.Quantity, &item.Instruction); err != nil {
			log.Printf("ChefView: Error scanning order item: %v", err)
			return nil, err
		}

		if order, ok := orderMap[item.OrderID]; ok {
			order.Items = append(order.Items, item)
		}
	}
	if err = itemRows.Err(); err != nil {
		return nil, err
	}

	finalOrders := make([]types.Order, 0, len(orderedList))
	for _, orderPtr := range orderedList {
		if len(orderPtr.Items) > 0 {
			finalOrders = append(finalOrders, *orderPtr)
		}
	}

	return finalOrders, nil
}
