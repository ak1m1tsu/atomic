package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Resposne struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Resposne {
	return Resposne{
		Status: StatusOK,
	}
}

func Error(msg string) Resposne {
	return Resposne{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationErrors) Resposne {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("%s is required", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("%s is not a valid url", err.Field()))
		case "alphanum":
			errMsgs = append(errMsgs, fmt.Sprintf("%s can only contain alphanumeric characters", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("%s is not valid", err.Field()))
		}
	}
	return Resposne{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
