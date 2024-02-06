package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"text/template"

	"github.com/google/uuid"

	"github.com/shutt90/email-notification/internal/params"
	mail "github.com/shutt90/email-notification/mailer"
)

const (
	PORT              = "9001"
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
	Href         string

	Recipients []string `json:"recipients"`
}

func main() {
	http.HandleFunc("/", HandleRequest)

	if err := http.ListenAndServe(":"+PORT, nil); err != nil {
		log.Fatal("unable to start server, error: ", err)
	}
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	params := params.New()
	params.SetEmail(r.URL.Query().Get("email"))
	parsedId, err := uuid.Parse(r.URL.Query().Get("uuid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request - missing required query params"))

		return
	}

	params.SetId(parsedId)

	cfg := cfg{}
	bodyByte, err := io.ReadAll(r.Body)
	if err := json.Unmarshal(bodyByte, &cfg); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("bad request %s", err.Error())))

		return
	}

	html, err := os.ReadFile("./template/emailBody.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unable to process at this time: %s", err.Error())))

		return
	}

	tmpl := template.Must(template.New("emailNotification").Parse(string(html)))

	cfg.Href = r.Host + "?" + params.GetId().String()

	buf := new(bytes.Buffer)
	tmpl.Execute(buf, cfg)

	mailClient := mail.New("", cfg.Username, cfg.Password, cfg.Address)
	msg := mail.BuildMessage(cfg.Recipients, fmt.Sprintf("%s:%d", cfg.Username, cfg.Port), VerifyUser, buf.String())
	if err := smtp.SendMail(fmt.Sprintf("%s:%d", cfg.Address, cfg.Port), mailClient.Auth, cfg.Username, cfg.Recipients, []byte(msg)); err != nil {
		if err.Error() == "EOF" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("user credentials incorrect"))

			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("server error logged %s", err.Error())))

		return
	}

	w.WriteHeader(http.StatusAccepted)
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
