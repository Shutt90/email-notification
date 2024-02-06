package mail

import (
	"net/smtp"
	"strings"
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

func BuildMessage(to []string, from, subject, message string) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n" +
		"From: " + from + "\n" +
		strings.Join(to, ";") + "\n" +
		"Subject: " + subject + "!\n" +
		message

	return msg
}
