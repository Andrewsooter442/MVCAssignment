package handler

import (
	"fmt"
	"net/http"

	"github.com/Andrewsooter442/MVCAssignment/config"
)

func (app *Application) HandleRootRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value(config.UserObject).(*config.JWTtoken)
	if !ok {
		http.Error(w, "Authentication error", http.StatusUnauthorized)
		return
	}

	client := config.Client{
		Name:    claims.Name,
		IsAdmin: claims.IsAdmin,
		IsChef:  claims.IsCheff,
		TableNo: 5,
	}
	if client.IsChef {
		orders, err := app.Pool.GetIncompleteOrdersForChef()
		if err != nil {
			http.Error(w, "Failed to load active orders", http.StatusInternalServerError)
			return
		}

		pageData := map[string]interface{}{
			"Client": client,
			"Orders": orders,
		}

		err = config.Templates.ExecuteTemplate(w, "chefHome.html", pageData)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to render chef page", http.StatusInternalServerError)
		}
		return
	}

	allCategory, err := app.Pool.GetAllCategories()
	if err != nil {
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	allItems, err := app.Pool.GetAllItems()
	if err != nil {
		http.Error(w, "Failed to load menu items", http.StatusInternalServerError)
		return
	}

	allPendingOrders, err := app.Pool.GetIncompleteOrders()
	if err != nil {
		http.Error(w, "Failed to load pending orders", http.StatusInternalServerError)
		return
	}

	menu := config.Menu{
		Categories: allCategory,
		Items:      allItems,
	}

	orderStatus := r.URL.Query().Get("status")

	homePageData := config.HomePageData{
		Client:        client,
		Menu:          menu,
		PendingOrders: allPendingOrders,
		StatusMessage: orderStatus,
	}

	err = config.Templates.ExecuteTemplate(w, "home.html", homePageData)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}
