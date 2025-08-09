package model

import (
	"database/sql"
	"errors"
	"github.com/Andrewsooter442/MVCAssignment/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

func (model *ModelConnection) CreateNewUser(user config.SignupRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `
        INSERT INTO users (name, mail, phone, password_hash, isAdmin, isCheff, created_at)
        VALUES (?, ?, ?, ?, FALSE, FALSE, UTC_TIMESTAMP())
    `
	_, err = model.DB.Exec(query, user.Username, user.Email, user.Phone, string(hashedPassword))
	if err != nil {
		log.Printf("Error executing insert statement for new user: %v", err)
		return errors.New("failed to create user in database")
	}

	//fmt.Println("User created successfully")
	return nil
}

func (model *ModelConnection) AuthenticateUser(loginData config.LoginRequest) (*config.JWTtoken, error) {
	var token config.JWTtoken
	var storedHash string

	query := `SELECT id, name, isAdmin, isCheff, password_hash FROM users WHERE name = ?`

	row := model.DB.QueryRow(query, loginData.Username)

	err := row.Scan(&token.ID, &token.Name, &token.IsAdmin, &token.IsCheff, &storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid username or password")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(loginData.Password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	//Add jwt voodoo
	token.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		Issuer:    "kitchen-app",
	}

	return &token, nil
}
