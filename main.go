package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/shutt90/email-notification/internal/params"
)

const (
	GoogleSMTPAddress = "smtp.gmail.com"
)

var (
	ErrUnknownAddress = fmt.Errorf("Error Unknown SMTP Address")
)

type cfg struct {
	Address      string `json:"smtpAddress"`
	Authenicated bool   `json:"authenticated"`
	Username     string `json:"username"`
	Password     string `json:"password"`
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
