package handler

import (
	"errors"
	"fmt"
	"log"
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
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
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

	data := config.LoginTemplate{
		InvalidCred:  false,
		ErrorMessage: "",
	}

	switch r.Method {
	case "POST":

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		var req config.LoginRequest
		req.Username = r.FormValue("username")
		req.Password = r.FormValue("password")

		if err := validateLoginRequest(req); err != nil {
			data.InvalidCred = true
			data.ErrorMessage = err.Error()
			err = config.Templates.ExecuteTemplate(w, "login.html", data)

			if err != nil {
				log.Printf("Error executing template: %v", err)
				http.Error(w, "There was a problem rendering the login page.", http.StatusInternalServerError)
			}
			return
		}

		fmt.Println("Login attempt for:", req.Username)

		// Create a JWT token
		fmt.Println("Querying the db for ", req.Username)
		token, err := app.Pool.AuthenticateUser(req)
		if err != nil {
			data.InvalidCred = true
			data.ErrorMessage = err.Error()
			err = config.Templates.ExecuteTemplate(w, "login.html", data)
			if err != nil {
				log.Printf("Error executing template: %v", err)
				http.Error(w, "There was a problem rendering the login page.", http.StatusInternalServerError)
			}
			return
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, token)
		signedString, err := jwtToken.SignedString([]byte(jwtSecret))
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to create new token", http.StatusBadRequest)
			return
		}

		//Send the cookie
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

		err := config.Templates.ExecuteTemplate(w, "login.html", data)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "There was a problem rendering the login page.", http.StatusInternalServerError)
		}
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}

func (app *Application) HandleSignupRequest(w http.ResponseWriter, r *http.Request) {

	data := config.LoginTemplate{
		InvalidCred:  false,
		ErrorMessage: "",
	}

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
			data.InvalidCred = true
			data.ErrorMessage = err.Error()
			err := config.Templates.ExecuteTemplate(w, "signup.html", data)
			if err != nil {
				log.Printf("Error executing template: %v", err)
				http.Error(w, "There was a problem rendering the login page.", http.StatusInternalServerError)
			}
			return

		}

		if err := app.Pool.CreateNewUser(req); err != nil {
			data.InvalidCred = true
			data.ErrorMessage = err.Error()
			err := config.Templates.ExecuteTemplate(w, "signup.html", data)
			if err != nil {
				log.Printf("Error executing template: %v", err)
				http.Error(w, "There was a problem rendering the login page.", http.StatusInternalServerError)
			}
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
		return

	case "GET":
		err := config.Templates.ExecuteTemplate(w, "signup.html", data)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "There was a problem rendering the login page.", http.StatusInternalServerError)
		}
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
