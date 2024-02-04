package smtp

import (
	"net/smtp"

	"github.com/shutt90/email-notification/internal/params"
)

type SMTP struct {
	Identity, Username, Password, Host string
	Auth                               smtp.Auth
}

func New(identity, username, password, host string) *SMTP {
	s := &SMTP{
		Identity: identity,
		Username: username,
		Password: password,
		Host:     host,
		Auth:     smtp.PlainAuth(identity, username, password, host),
	}

	return s
}

func (s *SMTP) Send(params.AcceptedQueryParams) {
	s.Auth.Start(smtp.ServerInfo{})
}
