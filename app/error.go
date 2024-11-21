package app

import (
	"encoding/json"
	"fmt"
)

var (
	// TechnicalErrors
	Err_Technical = AppError{Code: "2000", Message: "please contact tech support"}

	// BussinessErrors
	Err_BussinessErrors_1 = AppError{Code: "1001", Message: "not Found"}
	Err_BussinessErrors_2 = AppError{Code: "1004", Message: "invalid request"}

	// UnknownErrors
	Err_Unknown = AppError{Code: "9999", Message: "please contact admin"}
)

type AppError struct {
	Code    string
	Message string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("app error:%s|%s", e.Code, e.Message)
}

func (e *AppError) String() string {
	str, _ := json.Marshal(e)
	return string(str)
}

type ErrorResponse struct {
	Result ErrorResult `json:"result"`
}

type ErrorResult struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewErrorResponse(err error) ErrorResponse {
	var code, message string
	switch v := err.(type) {
	case *AppError:
		code = v.Code
		message = v.Message
	default:
		code = Err_Unknown.Code
		message = Err_Unknown.Message
	}

	return ErrorResponse{
		Result: ErrorResult{
			Code:    code,
			Message: message,
		},
	}
}
