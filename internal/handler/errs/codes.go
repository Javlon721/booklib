package errs

import "net/http"

const (
	NotFound ErrCode = iota
	AlreadyExists
	Aborted
	Internal
	InvalidArgument
	PermissionDenied
	Unauthenticated
)

var codeNames = map[ErrCode]string{
	NotFound:         "not_found",
	AlreadyExists:    "already_exists",
	Aborted:          "aborted",
	Internal:         "internal_server_error",
	InvalidArgument:  "invalid_argument",
	PermissionDenied: "permission_denied",
	Unauthenticated:  "unauthenticated",
}

var httpStatus = map[ErrCode]int{
	NotFound:         http.StatusNotFound,
	AlreadyExists:    http.StatusConflict,
	Aborted:          http.StatusConflict,
	Internal:         http.StatusInternalServerError,
	InvalidArgument:  http.StatusBadRequest,
	PermissionDenied: http.StatusForbidden,
	Unauthenticated:  http.StatusUnauthorized,
}
