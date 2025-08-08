package handler

import (
	"fmt"
	"github.com/Andrewsooter442/MVCAssignment/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
	"time"
)

func (app *Application) HandleLoginRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandleLoginRequest")
	switch r.Method {
	case "POST":
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
		fmt.Println(r.Form)

		var req model.LoginRequest
		req.Username = r.FormValue("username")
		req.Password = r.FormValue("password")

		fmt.Println("Login attempt for:", req.Username)
		token, err := app.Pool.AuthenticateUser(req)
		if err != nil {
			http.Error(w, "Failed to create new user", http.StatusBadRequest)
			return
		}

		fmt.Println("Login attempt for:", req.Username)
		fmt.Println(token)

		// Create a JWT token
		jwtSecret := os.Getenv("JWT_SECRET")
		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, token)
		signedString, err := jwtToken.SignedString([]byte(jwtSecret))
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to create new token", http.StatusBadRequest)
			return
		}

		//Send the cookee
		expirationTime := time.Now().Add(24 * time.Hour)
		cookie := http.Cookie{
			Name:    "token",
			Value:   signedString,
			Expires: expirationTime,

			HttpOnly: true,
			Secure:   false,
			Path:     "/",
		}

		fmt.Println("Cookie:", cookie)

		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/", http.StatusFound)
		return

	case "GET":
		fmt.Println("Serving login page")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "This is the login page.")
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}

func (app *Application) HandleSignupRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		var req model.SignupRequest
		req.Username = r.FormValue("username")
		req.Password = r.FormValue("password")
		req.Phone = r.FormValue("phone")
		req.Email = r.FormValue("email")

		if len(req.Password) < 6 {
			http.Error(w, "Password too short", http.StatusBadRequest)
			return
		}
		if !strings.Contains(req.Email, "@") {
			http.Error(w, "Invalid email format", http.StatusBadRequest)
			return
		}

		if err := app.Pool.CreateNewUser(req); err != nil {
			http.Error(w, "Failed to create new user", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Signup successful!")
		return

	case "GET":
		fmt.Println("Serving signup page")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "This is the signup page.")
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
