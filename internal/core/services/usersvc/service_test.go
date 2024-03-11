package usersvc

import (
	"context"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/shutt90/email-notification/internal/repositories/ports"
)

type mockDb struct {
	ctx  context.Context
	conn ports.UserRepo
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
			UPDATE user SET authenticated = TRUE WHERE uuid = $1 AND email = $2 AND authenticated = FALSE;
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
