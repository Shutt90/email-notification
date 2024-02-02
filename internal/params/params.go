package params

import "github.com/google/uuid"

type acceptedQueryParams struct {
	email string
	id    uuid.UUID
}

func New() *acceptedQueryParams {
	return &acceptedQueryParams{}
}

func (aqp *acceptedQueryParams) SetEmail(email string) {
	aqp.email = email
}

func (aqp *acceptedQueryParams) GetEmail() string {
	return aqp.email
}

func (aqp *acceptedQueryParams) SetId(id uuid.UUID) {
	aqp.id = id
}

func (aqp *acceptedQueryParams) GetId() uuid.UUID {
	return aqp.id
}
