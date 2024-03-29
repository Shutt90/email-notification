package configdbrepo

import (
	"context"
	"regexp"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"

	"github.com/shutt90/email-notification/internal/repositories/ports"
)

type mockDb struct {
	ctx  context.Context
	conn ports.UserRepo
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
		);`),
	).WillReturnResult(pgxmock.NewResult("CREATE", 0))

	mockConn.ExpectCommit()

	if err := mockClient.CreateTable("../../schema/user.sql"); err != nil {
		t.Fatal(err)
	}

	if err := mockConn.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}

}
