package model

import (
	"context"
	"fmt"
	"github.com/Andrewsooter442/MVCAssignment/config"
	"log"
	"time"
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
