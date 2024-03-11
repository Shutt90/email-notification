package usersvc

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/shutt90/email-notification/internal/core/domain"
	"github.com/shutt90/email-notification/internal/core/services/errors"
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
	_, err := svc.user.Exec(svc.db.GetContext(), "INSERT INTO user (email, uuid) VALUES ($1, $2)", user.Email, user.UUID)
	if err != nil {
		return err
	}

	return nil
}

func (svc *service) AuthenticateUser(id uuid.UUID, email string) error {
	tx, err := svc.user.Begin(svc.db.GetContext())
	if err != nil {
		return err
	}
	defer svc.user.Close(svc.db.GetContext())

	query := fmt.Sprintf("UPDATE user SET authenticated = TRUE WHERE uuid = $1 AND email = $2 AND authenticated = FALSE")

	tag, err := tx.Exec(svc.db.GetContext(), query, id, email)
	if err != nil {
		tx.Rollback(svc.db.GetContext())

		return err
	}

	if tag.RowsAffected() != 1 {
		tx.Rollback(svc.db.GetContext())

		return errors.ErrTooManyRows
	}

	if tag.RowsAffected() == 0 {
		tx.Rollback(svc.db.GetContext())

		return errors.ErrUpdateFailed
	}

	tx.Commit(svc.db.GetContext())

	return nil
}
