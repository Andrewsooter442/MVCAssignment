package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Andrewsooter442/MVCAssignment/config"
)

func (model *ModelConnection) CreateCategory(category *config.Category) error {
	query := `INSERT INTO categories (name) VALUES (?)`

	_, err := model.DB.Exec(query, category.Name)
	if err != nil {
		log.Printf("Error executing insert statement for new category: %v", err)
		return fmt.Errorf("failed to create category in database: %w", err)
	}

	return nil
}

func (model *ModelConnection) UpdateCategory(category *config.Category) error {
	query := `UPDATE categories SET name = ? WHERE id = ?`

	res, err := model.DB.Exec(query, category.Name, category.ID)
	if err != nil {
		log.Printf("Error executing update for category %d: %v", category.ID, err)
		return fmt.Errorf("failed to update category: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected for category update: %v", err)
		return fmt.Errorf("failed to verify category update: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no category found with ID %d to update", category.ID)
	}

	return nil
}

func (model *ModelConnection) CreateItem(item *config.Item) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO items (category_id, name, price, description) VALUES (?, ?, ?, ?)`
	_, err := model.DB.ExecContext(ctx, query, item.CategoryID, item.Name, item.Price, item.Description)
	if err != nil {
		log.Printf("Error executing insert for new item: %v", err)
		return fmt.Errorf("failed to create item: %w", err)
	}
	return nil
}

func (model *ModelConnection) UpdateItem(item *config.Item) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE items SET name = ?, price = ?, description = ?, category_id = ? WHERE id = ?`
	_, err := model.DB.ExecContext(ctx, query, item.Name, item.Price, item.Description, item.CategoryID, item.ID)
	if err != nil {
		log.Printf("Error executing update for item %d: %v", item.ID, err)
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

func (model *ModelConnection) GetCategoryByID(id int) (*config.Category, error) {
	query := `SELECT id, name FROM categories WHERE id = ?`
	var category config.Category
	err := model.DB.QueryRow(query, id).Scan(&category.ID, &category.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category with ID %d not found", id)
		}
		log.Printf("Error scanning category with ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return &category, nil
}

func (model *ModelConnection) GetAllCategories() ([]config.Category, error) {
	query := `SELECT id, name FROM categories ORDER BY name ASC`
	rows, err := model.DB.Query(query)
	if err != nil {
		log.Printf("Error querying all categories: %v", err)
		return nil, fmt.Errorf("failed to get all categories: %w", err)
	}
	defer rows.Close()

	var categories []config.Category
	for rows.Next() {
		var category config.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			log.Printf("Error scanning a category: %v", err)
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error during rows iteration for categories: %v", err)
		return nil, fmt.Errorf("error iterating category rows: %w", err)
	}

	return categories, nil
}
func (model *ModelConnection) GetAllItems() ([]config.Item, error) {
	query := `SELECT id, category_id, name, price, description FROM items`
	rows, err := model.DB.Query(query)
	if err != nil {
		log.Printf("Error querying items: %v", err)
		return nil, err
	}
	defer rows.Close()

	var items []config.Item
	for rows.Next() {
		var item config.Item
		if err := rows.Scan(&item.ID, &item.CategoryID, &item.Name, &item.Price, &item.Description); err != nil {
			log.Printf("Error scanning item: %v", err)
			return nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Rows iteration error on items: %v", err)
		return nil, err
	}

	return items, nil
}

func (model *ModelConnection) GetItemByID(id int) (*config.Item, error) {
	query := `SELECT id, category_id, name, price, description FROM items WHERE id = ?`
	var item config.Item
	err := model.DB.QueryRow(query, id).Scan(&item.ID, &item.CategoryID, &item.Name, &item.Price, &item.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("item with ID %d not found", id)
		}
		log.Printf("Error scanning item with ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to get item: %w", err)
	}
	return &item, nil
}

func (model *ModelConnection) GetAllOrders() ([]config.Order, error) {
	orderQuery := `SELECT id, user_id, table_no, complete, created_at FROM orders ORDER BY created_at DESC`
	orderRows, err := model.DB.Query(orderQuery)
	if err != nil {
		log.Printf("Error querying all orders: %v", err)
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer orderRows.Close()

	ordersMap := make(map[int]*config.Order)
	var ordersList []*config.Order

	for orderRows.Next() {
		var o config.Order
		if err := orderRows.Scan(&o.ID, &o.UserID, &o.TableNumber, &o.Complete, &o.CreatedAt); err != nil {
			log.Printf("Error scanning an order: %v", err)
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		o.Items = []config.OrderItem{}
		ordersMap[o.ID] = &o
		ordersList = append(ordersList, &o)
	}
	if err = orderRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating order rows: %w", err)
	}

	if len(ordersMap) == 0 {
		return []config.Order{}, nil
	}

	itemQuery := `SELECT order_id, item_id, quantity, instruction FROM order_items`
	itemRows, err := model.DB.Query(itemQuery)
	if err != nil {
		log.Printf("Error querying all order items: %v", err)
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer itemRows.Close()

	for itemRows.Next() {
		var oi config.OrderItem
		if err := itemRows.Scan(&oi.OrderID, &oi.ItemID, &oi.Quantity, &oi.Instruction); err != nil {
			log.Printf("Error scanning an order item: %v", err)
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}

		if order, ok := ordersMap[oi.OrderID]; ok {
			order.Items = append(order.Items, oi)
		}
	}
	if err = itemRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating order item rows: %w", err)
	}

	finalOrders := make([]config.Order, len(ordersList))
	for i, orderPtr := range ordersList {
		finalOrders[i] = *orderPtr
	}

	return finalOrders, nil
}

func (model *ModelConnection) GetIncompleteOrders() ([]config.Order, error) {
	orderQuery := `
		SELECT id, user_id, table_no, complete, created_at 
		FROM orders 
		WHERE complete = FALSE 
		ORDER BY created_at ASC`

	orderRows, err := model.DB.Query(orderQuery)
	if err != nil {
		log.Printf("Error querying incomplete orders: %v", err)
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer orderRows.Close()

	ordersMap := make(map[int]*config.Order)
	var ordersList []*config.Order

	for orderRows.Next() {
		var o config.Order
		if err := orderRows.Scan(&o.ID, &o.UserID, &o.TableNumber, &o.Complete, &o.CreatedAt); err != nil {
			log.Printf("Error scanning an order: %v", err)
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		o.Items = []config.OrderItem{}
		ordersMap[o.ID] = &o
		ordersList = append(ordersList, &o)
	}
	if err = orderRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating order rows: %w", err)
	}

	if len(ordersMap) == 0 {
		return []config.Order{}, nil
	}

	itemQuery := `SELECT order_id, item_id, quantity, instruction FROM order_items`
	itemRows, err := model.DB.Query(itemQuery)
	if err != nil {
		log.Printf("Error querying all order items: %v", err)
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer itemRows.Close()

	for itemRows.Next() {
		var oi config.OrderItem
		if err := itemRows.Scan(&oi.OrderID, &oi.ItemID, &oi.Quantity, &oi.Instruction); err != nil {
			log.Printf("Error scanning an order item: %v", err)
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}

		if order, ok := ordersMap[oi.OrderID]; ok {
			order.Items = append(order.Items, oi)
		}
	}
	if err = itemRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating order item rows: %w", err)
	}

	finalOrders := make([]config.Order, len(ordersList))
	for i, orderPtr := range ordersList {
		finalOrders[i] = *orderPtr
	}

	return finalOrders, nil
}
