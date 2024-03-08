package domain

import (
	"time"
)

type User struct {
	ID              int
	Email           string
	UUID            string
	Authenticated   bool
	AuthenticatedAt time.Time
}
