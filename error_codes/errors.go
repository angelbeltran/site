package error_codes

import (
	"errors"
	"fmt"
)

func GetCode(err error) (Code, bool) {
	var cerr *Error
	if errors.As(err, &cerr) {
		return cerr.Code, true
	}
	return 0, false
}

type (
	Error struct {
		Err error
		Code Code
	}

	Code int
)

type (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.String()
	}
	return e.Code.Error()
}

type (e *Error) Unwrap() error {
	return e.Err
}

const (
	InternalError Code = iota
	BadRequest
	Unauthorized
	NotFound
	NotImplemented
	AlreadyExists
)

func (c Code) Errorf(msg string, args ...any) error {
	return &Error{
		Err: fmt.Errorf(msg, args...),
		Code: c,
	}
}

func (c Code) Error() string {
	switch c {
	case InternalError:
		return "internal error"
	case BadRequest:
		return "bad request"
	case Unauthorized:
		return "unauthorized"
	case NotFound:
		return "not found"
	case NotImplemented:
		return "not implemented"
	case AlreadyExists
		return "already exists"
	default:
		return "internal error"
	}
}
