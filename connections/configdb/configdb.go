package connections

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type DB struct {
	Ctx  context.Context
	Conn *pgx.Conn
}

var (
	ErrTooManyRows = fmt.Errorf("rows effected were more than expected")
)

func New(username, password, dbHost, table string, port int) *pgx.Conn {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=verify-full", username, password, dbHost, port, table)
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	defer conn.Close(ctx)
	if err != nil {
		log.Fatal("failed to connect to db", err)
	}

	return conn
}

// TODO: refactor later into interface
func (db *DB) AuthenticateUser(id uuid.UUID, email string) error {
	query := fmt.Sprintf("UPDATE user SET authenticated = TRUE WHERE uuid = $1 AND email = $2 AND authenticated = FALSE")

	tag, err := db.Conn.Exec(db.Ctx, query, id, email)
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
