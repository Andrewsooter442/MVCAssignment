package config

import "github.com/golang-jwt/jwt/v5"

const UserObject = "user"

type JWTtoken struct {
	ID      int
	Name    string
	IsAdmin bool
	IsCheff bool
	jwt.RegisteredClaims
}
