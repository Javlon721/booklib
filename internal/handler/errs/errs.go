package errs

import (
	"encoding/json"
	"fmt"
	"runtime"
)

type ErrCode int

func (e ErrCode) String() string {
	return codeNames[e]
}

func (e ErrCode) Equal(e2 ErrCode) bool {
	return e == e2
}

type Error struct {
	Code     ErrCode
	Message  string
	FuncName string
	FileName string
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) HTTPStatus() int {
	return httpStatus[e.Code]
}

func New(code ErrCode, err error) *Error {
	pc, filename, line, _ := runtime.Caller(1)

	return &Error{
		Code:     code,
		Message:  err.Error(),
		FuncName: runtime.FuncForPC(pc).Name(),
		FileName: fmt.Sprintf("%s:%d", filename, line),
	}
}

func Errorf(code ErrCode, format string, v ...any) *Error {
	pc, filename, line, _ := runtime.Caller(1)

	return &Error{
		Code:     code,
		Message:  fmt.Sprintf(format, v...),
		FuncName: runtime.FuncForPC(pc).Name(),
		FileName: fmt.Sprintf("%s:%d", filename, line),
	}
}

// -------------------------------------------------------------------------

type FieldError struct {
	Field string `json:"field"`
	Err   string `json:"error"`
}

type FieldErrors []FieldError

func (fe *FieldErrors) Add(field string, err error) {
	*fe = append(*fe, FieldError{
		Field: field,
		Err:   err.Error(),
	})
}

func (fe *FieldErrors) ToError() *Error {
	return New(InvalidArgument, fe)
}

func (fe *FieldErrors) Error() string {
	str, err := json.Marshal(fe)

	if err != nil {
		return err.Error()
	}

	return string(str)
}
