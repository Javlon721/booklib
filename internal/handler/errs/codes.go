package errs

import "net/http"

const (
	NotFound ErrCode = iota
	AlreadyExists
	Aborted
	Internal
	InvalidArgument
)

var codeNames = map[ErrCode]string{
	NotFound:        "not_found",
	AlreadyExists:   "already_exists",
	Aborted:         "aborted",
	Internal:        "internal_server_error",
	InvalidArgument: "invalid_argument",
}

var httpStatus = map[ErrCode]int{
	NotFound:        http.StatusNotFound,
	AlreadyExists:   http.StatusConflict,
	Aborted:         http.StatusConflict,
	Internal:        http.StatusInternalServerError,
	InvalidArgument: http.StatusBadRequest,
}
