package ports

import (
	"github.com/shutt90/email-notification/internal/core/domain"
)

type UserService interface {
	CreateUser(user UserRepo) error
	GetUser(id string) (domain.User, error)
	AuthenticateUser(id string, email string) error
}
