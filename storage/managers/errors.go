package managers

import "errors"

var (
	ErrNotFound          = errors.New("managers: not found")
	ErrNotImplemented    = errors.New("managers: not found")
	ErrInvalidExpiration = errors.New("managers: invalid expiration")
)
