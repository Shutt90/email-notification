package configdbrepo

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shutt90/email-notification/internal/core/domain"
	"github.com/shutt90/email-notification/internal/repositories/ports"
)

type db struct {
	ctx  context.Context
	conn ports.UserRepo
}

var (
	ErrTooManyRows  = fmt.Errorf("rows effected were more than expected")
	ErrUpdateFailed = fmt.Errorf("update failed")
)

func New(username, password, dbHost, table string) *db {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=verify-full", username, password, dbHost, table)
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	defer conn.Close(ctx)
	if err != nil {
		log.Fatal("failed to connect to db", err)
	}

	return &db{
		ctx:  context.Background(),
		conn: conn,
	}
}

func (db *db) GetContext() context.Context {
	return db.ctx
}

func (db *db) CreateTable(filepath string) error {
	tx, err := db.conn.Begin(db.ctx)
	if err != nil {
		return err
	}
	defer db.conn.Close(db.ctx)

	f, err := os.ReadFile(filepath)
	if err != nil {
		tx.Rollback(db.ctx)

		return err
	}

	_, err = tx.Exec(db.ctx, string(f))
	if err != nil {
		tx.Rollback(db.ctx)

		return err
	}

	tx.Commit(db.ctx)

	return nil
}

func (db *db) CreateUser(user domain.User) error {
	_, err := db.conn.Exec(db.ctx, "INSERT INTO user (email, uuid) VALUES ($1, $2)", user.Email, user.UUID)
	if err != nil {
		return err
	}

	return nil
}

func (db *db) AuthenticateUser(id uuid.UUID, email string) error {
	tx, err := db.conn.Begin(db.ctx)
	if err != nil {
		return err
	}
	defer db.conn.Close(db.ctx)

	query := fmt.Sprintf("UPDATE user SET authenticated = TRUE WHERE uuid = $1 AND email = $2 AND authenticated = FALSE")

	tag, err := tx.Exec(db.ctx, query, id, email)
	if err != nil {
		tx.Rollback(db.ctx)

		return err
	}

	if tag.RowsAffected() != 1 {
		tx.Rollback(db.ctx)

		return ErrTooManyRows
	}

	if tag.RowsAffected() == 0 {
		tx.Rollback(db.ctx)

		return ErrUpdateFailed
	}

	tx.Commit(db.ctx)

	return nil
}
