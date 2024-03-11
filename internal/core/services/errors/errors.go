package errors

import "fmt"

var (
	ErrTooManyRows  = fmt.Errorf("error: too many rows effected")
	ErrUpdateFailed = fmt.Errorf("error: update failed")
)
