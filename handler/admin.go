package handler

import (
	"errors"
	"github.com/Andrewsooter442/MVCAssignment/config"
	"net/http"
	"strconv"
)

func (app *Application) HandleAdminRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/api/"):]
	switch path {
	case "addCategory":
		app.HandleAddCategory(w, r)
	case "editCategory":
		app.HandleEditCategory(w, r)
	case "addItem":
		app.HandleAddItem(w, r)
	case "editItem":
		app.HandleEditItem(w, r)

		//To implement
	//case "viewOlderOrders":

	default:
		http.NotFound(w, r)
	}
}

func getAdminClaims(r *http.Request) (*config.JWTtoken, error) {
	claims, ok := r.Context().Value(userObject).(*config.JWTtoken)
	if !ok {
		return nil, errors.New("could not retrieve user claims")
	}
	if !claims.IsAdmin {
		return nil, errors.New("user does not have admin privileges")
	}
	return claims, nil
}

func (app *Application) HandleAddCategory(w http.ResponseWriter, r *http.Request) {
	if _, err := getAdminClaims(r); err != nil {
		http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	categoryName := r.FormValue("name")
	if categoryName == "" {
		http.Error(w, "Category name cannot be empty", http.StatusBadRequest)
		return
	}

	category := config.Category{Name: categoryName}
	err := app.Pool.CreateCategory(&category)
	if err != nil {
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *Application) HandleEditCategory(w http.ResponseWriter, r *http.Request) {
	if _, err := getAdminClaims(r); err != nil {
		http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	categoryID, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	categoryName := r.FormValue("name")
	if categoryName == "" {
		http.Error(w, "Category name cannot be empty", http.StatusBadRequest)
		return
	}

	category := config.Category{ID: categoryID, Name: categoryName}
	if err := app.Pool.UpdateCategory(&category); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *Application) HandleAddItem(w http.ResponseWriter, r *http.Request) {
	if _, err := getAdminClaims(r); err != nil {
		http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

	categoryID, err := strconv.Atoi(r.FormValue("categoryId"))
	if err != nil {
		http.Error(w, "Invalid categoryId format", http.StatusBadRequest)
		return
	}

	item := config.Item{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Price:       price,
		CategoryID:  categoryID,
	}

	if item.Name == "" || item.Description == "" {
		http.Error(w, "Name and description cannot be empty", http.StatusBadRequest)
		return
	}

	err = app.Pool.CreateItem(&item)
	if err != nil {
		http.Error(w, "Failed to create item", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (app *Application) HandleEditItem(w http.ResponseWriter, r *http.Request) {
	if _, err := getAdminClaims(r); err != nil {
		http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	itemID, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		http.Error(w, "Invalid price format", http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.Atoi(r.FormValue("categoryId"))
	if err != nil {
		http.Error(w, "Invalid categoryId format", http.StatusBadRequest)
		return
	}

	item := config.Item{
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
		http.Error(w, "Failed to update item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
