package params

import "github.com/google/uuid"

type AcceptedQueryParams struct {
	email string
	id    uuid.UUID
}

func New() *AcceptedQueryParams {
	return &AcceptedQueryParams{}
}

func (aqp *AcceptedQueryParams) SetEmail(email string) {
	aqp.email = email
}

func (aqp *AcceptedQueryParams) GetEmail() string {
	return aqp.email
}

func (aqp *AcceptedQueryParams) SetId(id uuid.UUID) {
	aqp.id = id
}

func (aqp *AcceptedQueryParams) GetId() uuid.UUID {
	return aqp.id
}
