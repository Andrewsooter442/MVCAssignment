package model

import (
	"database/sql"
)

type ModelConnection struct {
	DB *sql.DB
}
