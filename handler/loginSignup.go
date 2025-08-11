package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Andrewsooter442/MVCAssignment/config"
)

func validateLoginRequest(req config.LoginRequest) error {
	if req.Username == "" {
		return errors.New("username is a required field")
	}
	if req.Password == "" {
		return errors.New("password is a required field")
	}
	return nil
}

func validateSignupRequest(req config.SignupRequest) error {
	if req.Username == "" {
		return errors.New("username is a required field")
	}
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		return errors.New("invalid email format")
	}
	if len(req.Phone) < 10 {
		return errors.New("phone number seems too short")
	}

	return nil
}

func (app *Application) HandleLoginRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandleLoginRequest")
	switch r.Method {
	case "POST":
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
		//	fmt.Println(r.Form)

		var req config.LoginRequest
		req.Username = r.FormValue("username")
		req.Password = r.FormValue("password")

		if err := validateLoginRequest(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println("Login attempt for:", req.Username)
		token, err := app.Pool.AuthenticateUser(req)
		if err != nil {
			http.Error(w, "Failed to login", http.StatusBadRequest)
			return
		}

		//fmt.Println(token)

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

		var req config.SignupRequest
		req.Username = r.FormValue("username")
		req.Password = r.FormValue("password")
		req.Phone = r.FormValue("phone")
		req.Email = r.FormValue("email")

		if err := validateSignupRequest(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
