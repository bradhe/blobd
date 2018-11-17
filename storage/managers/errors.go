package managers

import "errors"

var (
	ErrNotFound          = errors.New("managers: not found")
	ErrNotImplemented    = errors.New("managers: not implemented")
	ErrInvalidExpiration = errors.New("managers: invalid expiration")
)
