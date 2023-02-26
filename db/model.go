package db

import (
	"database/sql"
	"time"
)

type Model struct {
	ID int64 `json:"id" db:"id"`
	ModelLite
}

type ModelLite struct {
	CreatedAt  time.Time    `json:"createdAt" db:"created_at"`
	ModifiedAt time.Time    `json:"modifiedAt" db:"modified_at"`
	DeletedAt  sql.NullTime `json:"deletedAt" db:"deleted_at"`
}
