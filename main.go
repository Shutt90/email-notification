package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const (
	GoogleSMTPAddress = "smtp.gmail.com"
)

var (
	ErrUnknownAddress = fmt.Errorf("Error Unknown SMTP Address")
)

type acceptedQueryParams struct {
	email string
	id    uuid.UUID
}

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
	params := &acceptedQueryParams{}
	params.setEmail(r.URL.Query().Get("email"))
	parsedId, err := uuid.Parse(r.URL.Query().Get("uuid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))

		return
	}

	params.setId(parsedId)
}

func (aqp *acceptedQueryParams) setEmail(email string) {
	aqp.email = email
}

func (aqp *acceptedQueryParams) getEmail() string {
	return aqp.email
}

func (aqp *acceptedQueryParams) setId(id uuid.UUID) {
	aqp.id = id
}

func (aqp *acceptedQueryParams) getId() uuid.UUID {
	return aqp.id
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
