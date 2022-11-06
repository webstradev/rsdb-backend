package utils

import (
	"github.com/webstradev/rsdb-backend/auth"
	"github.com/webstradev/rsdb-backend/db"
)

type Environment struct {
	DB  *db.Database
	JWT *auth.JWTService
}
