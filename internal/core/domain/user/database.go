package usersvc

import (
	"fmt"

	"github.com/google/uuid"
)

type service struct {
	user ports.UserRepo
}

func New(user ports.UserRepo) *service {
	return &service{
		user: user,
	}
}

func (svc *service) CreateUser(user ports.UserRepo) error {
	svc.user.Exec(svc.user.ctx, "INSERT INTO user (email, uuid) VALUES ($1, $2)", user.Email, user.U)
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
