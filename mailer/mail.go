package mail

import (
	"fmt"
	"net/smtp"
)

type Mail struct {
	Identity, Username, Password, Host string
	Auth                               smtp.Auth
}

func New(identity, username, password, host string) *Mail {
	s := &Mail{
		Identity: identity,
		Username: username,
		Password: password,
		Host:     host,
		Auth:     smtp.PlainAuth(identity, username, password, host),
	}

	return s
}

func BuildMessage(to []string, from, subject, body string) string {
	toString := "To: "
	for _, t := range to {
		toString += t + "; "
	}

	msg := fmt.Sprintf(`
		From: %s \n"
		To: %s\n
		Subject: %s \n\n
		%s`,
		from, toString, subject, body,
	)

	return msg
}
