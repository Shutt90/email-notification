package configdb

import (
	"context"
	"regexp"
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
	mockConn, err := pgxmock.NewConn()
	if err != nil {
		t.Fatal(err)
	}
	defer mockConn.Close(context.Background())

	mockClient := &db{
		context.Background(),
		mockConn,
	}

	mockConn.ExpectBeginTx(pgx.TxOptions{})

	mockConn.ExpectExec(regexp.QuoteMeta(`
		CREATE TABLE IF NOT EXISTS user (
		    id SERIAL PRIMARY KEY,
		    email VARCHAR(255) NOT NULL DEFAULT '',
		    uuid VARCHAR(255) NOT NULL DEFAULT '',
		    authenticated BOOLEAN NOT NULL DEFAULT false,
		    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		    authenticated_at TIMESTAMPTZ DEFAULT NULL
		)`),
	).WillReturnResult(pgxmock.NewResult("CREATE TABLE", 0))

	if err := mockClient.CreateTable(); err != nil {
		t.Fatal(err)
	}

}
