package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

type BaseResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func OK() BaseResponse {
	return BaseResponse{
		Status: StatusOK,
	}
}

func Error(msg string) BaseResponse {
	return BaseResponse{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationErrors) BaseResponse {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid", err.Field()))
		}
	}

	return BaseResponse{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
