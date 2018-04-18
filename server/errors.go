package server

import "net/http"
import "fmt"
import "errors"

var (
	// Indicates that there was an application error during the
	// request and the request should be abandoned.
	ErrAbandonRequest = errors.New("server: abandon request")

	ErrInvalidUUID          = errors.New("server: invalid uuid")
	ErrMissingDecryptionKey = errors.New("server: missing decryption key")
)

type Error struct {
	Code    string `json:"code"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

var StandardErrors = map[string]Error{
	"unknown": {
		Status:  http.StatusInternalServerError,
		Message: "an unknown error occured",
	},
	"unauthorized": {
		Status:  http.StatusUnauthorized,
		Message: "authorization failed",
	},
	"internal_server_error": {
		Status:  http.StatusInternalServerError,
		Message: "an internal server error occured",
	},
	"not_implemented": {
		Status:  http.StatusNotImplemented,
		Message: "support for path `%s %s` is not yet implemented",
	},
	"not_found": {
		Status:  http.StatusNotFound,
		Message: "path `%s %s` was not found",
	},
	"invalid_query": {
		Status:  http.StatusPreconditionFailed,
		Message: "period stop or stop missing",
	},
	"method_not_allowed": {
		Status:  http.StatusNotFound,
		Message: "`%s` method is not allowed for path %s",
	},
	"invalid_password": {
		Status:  http.StatusPreconditionFailed,
		Message: "invalid password",
	},
	"invalid_body": {
		Status:  http.StatusBadRequest,
		Message: "invalid body",
	},
}

func GetError(name string, args ...interface{}) Error {
	if err, ok := StandardErrors[name]; ok {
		cp := *(&err)
		cp.Message = fmt.Sprintf(cp.Message, args...)
		cp.Code = name

		return cp
	}

	// NOTE: We should guarantee that htis is always here...
	return GetError("unknown")
}
