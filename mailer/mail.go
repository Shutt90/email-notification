package mail

import (
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

	msg := "From: " + from + "\n" +
		toString + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	return msg
}
