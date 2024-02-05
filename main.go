package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/smtp"

	"github.com/google/uuid"

	"github.com/shutt90/email-notification/internal/params"
	mail "github.com/shutt90/email-notification/mailer"
)

const (
	GoogleSMTPAddress = "smtp.gmail.com"

	VerifyUser = "Verify your Account Creation"
)

var (
	ErrUnknownAddress = fmt.Errorf("Error Unknown SMTP Address")
)

type cfg struct {
	Address      string `json:"smtpAddress"`
	Authenicated bool   `json:"authenticated"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Port         uint   `json:"port"`

	Recipients []string `json:"recipients"`
}

func main() {
	http.HandleFunc("/", HandleRequest)
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	params := params.New()
	params.SetEmail(r.URL.Query().Get("email"))
	parsedId, err := uuid.Parse(r.URL.Query().Get("uuid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))

		return
	}

	params.SetId(parsedId)

	cfg := cfg{}
	bodyByte, err := io.ReadAll(r.Body)
	if err := json.Unmarshal(bodyByte, &cfg); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))

		return
	}

	mailClient := mail.New("", cfg.Username, cfg.Password, cfg.Address)
	msg := mail.BuildMessage(cfg.Recipients, fmt.Sprintf("%s:%d", cfg.Username, cfg.Port), VerifyUser, "verify your account creation by clicking the link")
	if err := smtp.SendMail(fmt.Sprintf("%s:%d", cfg.Address, cfg.Port), mailClient.Auth, cfg.Username, cfg.Recipients, []byte(msg)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error logged"))

		return
	}
}

func (cfg *cfg) smtpAddress() error {
	switch cfg.Address {
	case GoogleSMTPAddress:
		cfg.Authenicated = true
	default:
		return ErrUnknownAddress
	}

	return nil

}
