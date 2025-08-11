package config

import "github.com/golang-jwt/jwt/v5"

type contextKey string

const UserObject contextKey = "user"

type JWTtoken struct {
	ID      int
	Name    string
	IsAdmin bool
	IsCheff bool
	jwt.RegisteredClaims
}
