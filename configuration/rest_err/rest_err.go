package rest_err

import "net/http"

type RestErr struct {
	Message string   `json:"message" example:"error trying to process request"`
	Err     string   `json:"error" example:"internal_server_error"`
	Code    int      `json:"code" example:"500"`
	Causes  []Causes `json:"causes"`
}

type Causes struct {
	Field   string `json:"field" example:"name"`
	Message string `json:"message" example:"name is required"`
}

func (r *RestErr) Error() string {
	return r.Message
}

func NewError(message, err string) *RestErr {
	switch err {
	case "bad_request":
		return NewBadRequestError(message)
	case "not_found":
		return NewNotFoundError(message)
	default:
		return NewInternalServerError(message)
	}
}

func NewBadRequestError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "bad_request",
		Code:    http.StatusBadRequest,
	}
}

func NewBadRequestValidationError(message string, causes []Causes) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "bad_request",
		Code:    http.StatusBadRequest,
		Causes:  causes,
	}
}

func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "internal_server_error",
		Code:    http.StatusInternalServerError,
	}
}

func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "not_found",
		Code:    http.StatusNotFound,
	}
}
