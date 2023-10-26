package error

import (
	"net/http"

	"github.com/rahmatrdn/go-skeleton/entity"
)

const (
	SUCCESS_CODE         = "00"
	SUCCESS_MSG          = "Success"
	INVALID_AUTH_CODE    = "01"
	INVALID_AUTH_MSG     = "Invalid Email or Password"
	INVALID_PAYLOAD_CODE = "02"
	INVALID_PAYLOAD_MSG  = "Invalid Payload Request Data"
	INVALID_TOKEN_CODE   = "05"
	INVALID_TOKEN_MSG    = "Invalid Access Token"
	BAD_REQUEST_CODE     = "30"
	BAD_REQUEST_MSG      = "Bad Request"

	GENERAL_ERROR_MESSAGE = "Something went wrong."
)

func ErrRecordNotFound() CustomErrorResponse {
	return CustomErrorResponse{
		Message:  "Data not found",
		ErrCode:  BAD_REQUEST_MSG,
		HTTPCode: http.StatusNotFound,
	}
}

func ErrUserNotFound() CustomErrorResponse {
	return CustomErrorResponse{
		Message:  "User not found",
		ErrCode:  BAD_REQUEST_MSG,
		HTTPCode: http.StatusNotFound,
	}
}

func ErrInvalidEmailOrPassword() CustomErrorResponse {
	return CustomErrorResponse{
		Message:  INVALID_AUTH_MSG,
		ErrCode:  INVALID_AUTH_CODE,
		HTTPCode: http.StatusUnauthorized,
	}
}

func ErrInvalidToken() CustomErrorResponse {
	return CustomErrorResponse{
		Message:  INVALID_TOKEN_MSG,
		ErrCode:  INVALID_TOKEN_CODE,
		HTTPCode: http.StatusUnauthorized,
	}
}

func ErrInvalidPayload(meta []entity.ErrorResponse) CustomErrorResponseWithMeta {
	return CustomErrorResponseWithMeta{
		Message:  INVALID_PAYLOAD_MSG,
		ErrCode:  INVALID_PAYLOAD_CODE,
		HTTPCode: http.StatusUnprocessableEntity,
		Meta:     meta,
	}
}

type CustomErrorResponse struct {
	Message  string `json:"message,omitempty"`
	ErrCode  string `json:"code,omitempty"`
	HTTPCode int    `json:"http_code"`
}
type CustomErrorResponseWithMeta struct {
	Message  string                 `json:"message,omitempty"`
	ErrCode  string                 `json:"code,omitempty"`
	HTTPCode int                    `json:"http_code"`
	Meta     []entity.ErrorResponse `json:"meta,omitempty"`
}

// Error is a function to convert error to string.
// It exists to satisfy error interface
func (c CustomErrorResponse) Error() string {
	return c.Message
}

func ErrGeneralInvalid() CustomErrorResponse {
	return CustomErrorResponse{
		Message:  GENERAL_ERROR_MESSAGE,
		ErrCode:  BAD_REQUEST_MSG,
		HTTPCode: http.StatusUnprocessableEntity,
	}
}

func ErrInvalidRequest() CustomErrorResponse {
	return CustomErrorResponse{
		Message:  INVALID_PAYLOAD_MSG,
		ErrCode:  BAD_REQUEST_MSG,
		HTTPCode: http.StatusUnprocessableEntity,
	}
}

func CustomError(message string, errCode string, httpCode int) CustomErrorResponse {
	return CustomErrorResponse{
		Message:  message,
		ErrCode:  errCode,
		HTTPCode: httpCode,
	}
}
