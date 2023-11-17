package error

import (
	"net/http"

	"github.com/rahmatrdn/go-skeleton/entity"
)

func ErrRecordNotFound() CustomErrorResponse {
	return CustomErrorResponse{
		Message:  entity.DATA_NOT_FOUND_MSG,
		ErrCode:  entity.BAD_REQUEST_MSG,
		HTTPCode: http.StatusNotFound,
	}
}

func ErrUserNotFound() CustomErrorResponse {
	return CustomErrorResponse{
		Message:  entity.USER_NOT_FOUND_MSG,
		ErrCode:  entity.BAD_REQUEST_MSG,
		HTTPCode: http.StatusNotFound,
	}
}

func ErrInvalidEmailOrPassword() CustomErrorResponse {
	return CustomErrorResponse{
		Message:  entity.INVALID_AUTH_MSG,
		ErrCode:  entity.INVALID_AUTH_CODE,
		HTTPCode: http.StatusUnauthorized,
	}
}

func ErrInvalidToken() CustomErrorResponse {
	return CustomErrorResponse{
		Message:  entity.INVALID_TOKEN_MSG,
		ErrCode:  entity.INVALID_TOKEN_CODE,
		HTTPCode: http.StatusUnauthorized,
	}
}

func ErrInvalidPayload(meta []entity.ErrorResponse) CustomErrorResponseWithMeta {
	return CustomErrorResponseWithMeta{
		Message:  entity.INVALID_PAYLOAD_MSG,
		ErrCode:  entity.INVALID_PAYLOAD_CODE,
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
		Message:  entity.GENERAL_ERROR_MESSAGE,
		ErrCode:  entity.BAD_REQUEST_MSG,
		HTTPCode: http.StatusUnprocessableEntity,
	}
}

func ErrInvalidRequest() CustomErrorResponse {
	return CustomErrorResponse{
		Message:  entity.INVALID_PAYLOAD_MSG,
		ErrCode:  entity.BAD_REQUEST_MSG,
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
