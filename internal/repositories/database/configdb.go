package configdbrepo

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"

	"github.com/shutt90/email-notification/internal/repositories/ports"
)

type db struct {
	ctx  context.Context
	conn ports.UserRepo
}

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
