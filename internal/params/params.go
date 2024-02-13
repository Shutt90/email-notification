package params

import "github.com/google/uuid"

type AcceptedQueryParams struct {
	email string
	id    uuid.UUID
}

func New(email string, id uuid.UUID) *AcceptedQueryParams {
	return &AcceptedQueryParams{}
}

func (aqp *AcceptedQueryParams) GetEmail() string {
	return aqp.email
}

func (aqp *AcceptedQueryParams) GetId() uuid.UUID {
	return aqp.id
}
