package errs

import (
	"fmt"
)

var (
	ErrNoRow = fmt.Errorf("sql: no rows")
	ErrBadRequest = fmt.Errorf("bad request")
	ErrUnauthorized = fmt.Errorf("unauthorized")
	ErrForbidden = fmt.Errorf("forbidden")
	ErrInternalServerError = fmt.Errorf("internal server error")
)
