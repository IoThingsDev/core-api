package helpers

import (
	"fmt"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%v: %v", e.Code, e.Message)
}

func ErrorWithCode(code string, message string) error {
	return Error{Code: code, Message: message}
}
