package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Andrewsooter442/MVCAssignment/types"
)

// Category Handlers
func (app *Application) HandleGetAddCategory(w http.ResponseWriter, r *http.Request) {
	data := types.MenuEditPageData{
		Title:  "Add New Category",
		Action: "/admin/addCategory",
		Type:   "category",
	}

	err := types.Templates.ExecuteTemplate(w, "menuEdit.html", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *Application) HandlePostAddCategory(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	categoryName := r.FormValue("name")
	if categoryName == "" {
		http.Error(w, "Category name cannot be empty", http.StatusBadRequest)
		return
	}

	category := types.Category{Name: categoryName}
	if err := app.Pool.CreateCategory(&category); err != nil {
		log.Printf("Failed to create category: %v", err)
		http.Error(w, "Failed to create category. It might already exist.", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Item Handlers
func (app *Application) HandleGetAddItem(w http.ResponseWriter, r *http.Request) {
	categories, err := app.Pool.GetAllCategories()
	if err != nil {
		log.Printf("Failed to get all categories: %v", err)
		http.Error(w, "Failed to load data for form", http.StatusInternalServerError)
		return
	}

	data := types.MenuEditPageData{
		Title:      "Add New Item",
		Action:     "/admin/addItem",
		Type:       "item",
		Categories: categories,
	}

	if err := types.Templates.ExecuteTemplate(w, "menuEdit.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *Application) HandlePostAddItem(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		http.Error(w, "Invalid price format", http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.Atoi(r.FormValue("category_id"))
	if err != nil {
		http.Error(w, "Invalid category selected", http.StatusBadRequest)
		return
	}

	item := types.Item{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Price:       price,
		CategoryID:  categoryID,
	}

	if item.Name == "" || item.Description == "" {
		http.Error(w, "Name and description cannot be empty", http.StatusBadRequest)
		return
	}

	if err := app.Pool.CreateItem(&item); err != nil {
		log.Printf("Failed to create item: %v", err)
		http.Error(w, "Failed to create item", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) HandleGetEditItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid Item ID", http.StatusBadRequest)
		return
	}

	item, err := app.Pool.GetItemByID(id)
	if err != nil {
		log.Printf("Failed to get item %d: %v", id, err)
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	categories, err := app.Pool.GetAllCategories()
	if err != nil {
		log.Printf("Failed to get all categories: %v", err)
		http.Error(w, "Failed to load data for form", http.StatusInternalServerError)
		return
	}

	data := types.MenuEditPageData{
		Title:      "Edit Item",
		Action:     fmt.Sprintf("/admin/editItem/%d", id),
		Type:       "item",
		Item:       *item,
		Categories: categories,
	}

	if err := types.Templates.ExecuteTemplate(w, "menuEdit.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *Application) HandlePostEditItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid Item ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		http.Error(w, "Invalid price format", http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.Atoi(r.FormValue("category_id"))
	if err != nil {
		http.Error(w, "Invalid category selected", http.StatusBadRequest)
		return
	}

	item := types.Item{
		ID:          itemID,
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Price:       price,
		CategoryID:  categoryID,
	}

	if item.Name == "" || item.Description == "" {
		http.Error(w, "Name and description cannot be empty", http.StatusBadRequest)
		return
	}

	if err := app.Pool.UpdateItem(&item); err != nil {
		log.Printf("Failed to update item %d: %v", itemID, err)
		http.Error(w, "Failed to update item", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) HandleGetViewOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid Order ID", http.StatusBadRequest)
		return
	}

	order, err := app.Pool.GetOrderById(id)
	if err != nil {
		log.Printf("Failed to get order %d: %v", id, err)
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	userName := app.Pool.GetUserNameById(order.UserID)

	var total float64
	for _, item := range order.Items {
		total += item.Price * float64(item.Quantity)
	}

	data := types.ViewOrderPageData{
		Order:    order,
		UserName: userName,
		Total:    total,
	}

	err = types.Templates.ExecuteTemplate(w, "orderView.html", data)
	if err != nil {
		log.Printf("Error executing order view template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *Application) HandleGetViewOldOrder(w http.ResponseWriter, r *http.Request) {
	orders, err := app.Pool.GetAllOrders()
	if err != nil {
		log.Printf("Failed to retrieve all orders: %v", err)
		http.Error(w, "Could not load order history.", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Orders": orders,
	}

	err = types.Templates.ExecuteTemplate(w, "viewOldOrders.html", data)
	if err != nil {
		log.Printf("Error executing order history template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
