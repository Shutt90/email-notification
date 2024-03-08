package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              int
	Email           string
	UUID            string
	Authenticated   bool
	AuthenticatedAt time.Time
}
