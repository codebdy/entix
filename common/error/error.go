package error

import "fmt"

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func New(code, message string) Error {
	return Error{
		Code:    code,
		Message: message,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("Error [%s] %s", e.Code, e.Message)
}
