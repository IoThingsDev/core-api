package helpers

import (
	"fmt"
)

type Error struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	HttpCode int    `json:"-"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%v: %v", e.Code, e.Message)
}

func ErrorWithCode(code string, message string) Error {
	return Error{Code: code, Message: message}
}

func NewError(httpCode int, code string, message string) Error {
	return Error{Code: code, Message: message, HttpCode: httpCode}
}
