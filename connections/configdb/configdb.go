package configdb

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type db struct {
	ctx  context.Context
	conn *pgx.Conn
}

var (
	ErrTooManyRows = fmt.Errorf("rows effected were more than expected")
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

// TODO: refactor later into interface
func (db *db) AuthenticateUser(id uuid.UUID, email string) error {
	query := fmt.Sprintf("UPDATE user SET authenticated = TRUE WHERE uuid = $1 AND email = $2 AND authenticated = FALSE")

	tag, err := db.conn.Exec(db.ctx, query, id, email)
	if err != nil {
		// TODO: setup logger
		fmt.Println(err)

		return err
	}

	if tag.RowsAffected() != 1 {
		return ErrTooManyRows
	}

	return nil
}
