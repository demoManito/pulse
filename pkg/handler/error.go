package handler

import (
	"errors"
	"fmt"
)

// ActionError is a status error.
type ActionError struct {
	Status  int
	Code    int
	Message string
}

// ActionError implements the error interface.
func (e *ActionError) Error() string {
	return fmt.Sprintf("error: status = %d code = %d message = %s", e.Status, e.Code, e.Message)
}

// Is matches each error in the chain with the target value.
func (e *ActionError) Is(err error) bool {
	if se := new(ActionError); errors.As(err, &se) {
		return se.Status == e.Status && se.Code == e.Code && se.Message == e.Message
	}
	return false
}

// NewActionError returns a new ActionError
func NewActionError(status int, code int, msg string) *ActionError {
	return &ActionError{
		Status:  status,
		Code:    code,
		Message: msg,
	}
}

// DecodeError decodes an error into a status code and an ActionError struct.
func DecodeError(err error) ActionError {
	switch err.(type) {
	case *ActionError:
		return *err.(*ActionError)
	case nil:
		return ActionError{Status: 200, Code: 0, Message: "成功"}
	case error:
		return ActionError{Status: 500, Code: 500, Message: err.Error()}
	default:
		return ActionError{Status: 500, Code: 500, Message: "未知错误"}
	}
}
