package usersvc

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/shutt90/email-notification/internal/core/domain"
	"github.com/shutt90/email-notification/internal/repositories/ports"
)

type service struct {
	user ports.UserRepo
	db   ports.UserService
}

func New(user ports.UserRepo, db ports.UserService) *service {
	return &service{
		user: user,
		db:   db,
	}
}

func (svc *service) CreateUser(user domain.User) error {
	svc.user.Exec(svc.db.GetContext(), "INSERT INTO user (email, uuid) VALUES ($1, $2)", user.Email, user.UUID)
}

func (svc *service) AuthenticateUser(id uuid.UUID, email string) error {
	tx, err := svc.user.Begin(svc.user.ctx)
	if err != nil {
		return err
	}
	defer svc.user.Close(svc.user.ctx)

	query := fmt.Sprintf("UPDATE user SET authenticated = TRUE WHERE uuid = $1 AND email = $2 AND authenticated = FALSE")

	tag, err := tx.Exec(svc.user.ctx, query, id, email)
	if err != nil {
		tx.Rollback(svc.user.ctx)

		return err
	}

	if tag.RowsAffected() != 1 {
		tx.Rollback(svc.user.ctx)

		return ErrTooManyRows
	}

	if tag.RowsAffected() == 0 {
		tx.Rollback(svc.user.ctx)

		return ErrUpdateFailed
	}

	tx.Commit(svc.user.ctx)

	return nil
}
