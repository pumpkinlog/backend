package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrConflict   = errors.New("already exists")
	ErrValidation = errors.New("validation error")
)

func ValidationError(msg string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrValidation, fmt.Sprintf(msg, args...))

}
