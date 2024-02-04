package connections

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
)

type DB *pgx.Conn

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
