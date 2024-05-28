// Вспомогательные функции при обработке запросов и ответов в хэндлере.
package api

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type RespStatus struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() RespStatus {
	return RespStatus{
		Status: StatusOK,
	}
}

func Error(msg string) RespStatus {
	return RespStatus{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationErrors) RespStatus {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "email":
			errMsgs = append(errMsgs, fmt.Sprintf("input %s is not a valid Email", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return RespStatus{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
