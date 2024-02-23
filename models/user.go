package models

import (
	"database/sql"

	"github.com/google/uuid"
)

type User struct {
	id              int
	email           string
	u               uuid.UUID
	authenticated   bool
	authenticatedAt sql.NullTime
}

func NewUser(id int, email string, u uuid.UUID) User {
	return User{
		id:    id,
		email: email,
		u:     u,
	}
}
