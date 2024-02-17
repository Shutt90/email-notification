package configdb

import (
	"context"
	"regexp"
	"testing"

	"github.com/google/uuid"
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
	).WillReturnResult(pgxmock.NewResult("CREATE", 0))

	mockConn.ExpectCommit()

	if err := mockClient.CreateTable(); err != nil {
		t.Fatal(err)
	}

	if err := mockConn.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}

}

func TestAuthenticateUser(t *testing.T) {
	mockConn, err := pgxmock.NewConn()
	if err != nil {
		t.Fatal(err)
	}
	defer mockConn.Close(context.Background())

	mockClient := &db{
		context.Background(),
		mockConn,
	}

	id := uuid.New()
	email := "test@example.com"

	mockConn.ExpectBeginTx(pgx.TxOptions{})

	mockConn.ExpectExec(regexp.QuoteMeta(`
			UPDATE user SET authenticated = TRUE WHERE uuid = $1 AND email = $2 AND authenticated = FALSE
		`),
	).WithArgs(id, email).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	mockConn.ExpectCommit()

	if err := mockClient.AuthenticateUser(id, email); err != nil {
		t.Fatal(err)
	}

	if err := mockConn.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
