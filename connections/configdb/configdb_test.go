package configdb

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v3"
)

type mockDb struct {
	ctx  context.Context
	conn PgxConnectionIface
}

func (db *mockDb) Begin(ctx context.Context) (pgx.Tx, error) {
	return db.conn.Begin(ctx)
}

func (db *mockDb) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return db.conn.Exec(ctx, sql, arguments)
}

func (db *mockDb) Close(ctx context.Context) error {
	return db.conn.Close(ctx)
}

func TestCreateTable(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mock.ExpectExec(`
		CREATE TABLE IF NOT EXISTS user (
		    id SERIAL PRIMARY KEY,
		    email VARCHAR(255) NOT NULL DEFAULT '',
		    uuid VARCHAR(255) NOT NULL DEFAULT '',
		    authenticated BOOLEAN NOT NULL DEFAULT false,
		    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		    authenticated_at TIMESTAMPTZ DEFAULT NULL
		)`,
	)

	mockConn, _ := pgxmock.NewConn()

	mockClient := &db{
		context.Background(),
		mockConn,
	}

	if err := mockClient.CreateTable(); err != nil {
		t.Fatal(err)
	}

}
