package error

import (
	"net/http"

	"gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/entity"
)

const (
	PathNotFoundCode = 997

	SUCCESS_CODE                  = "00"
	SUCCESS_MSG                   = "Success"
	INVALID_AUTH_CODE             = "01"
	INVALID_AUTH_MSG              = "Invalid Email or Password"
	INVALID_PAYLOAD_CODE          = "02"
	INVALID_PAYLOAD_MSG           = "Invalid Payload"
	INVALID_MERCHANT_CODE         = "03"
	INVALID_MERCHANT_MSG          = "Merchant not found"
	INVALID_TID_CODE              = "04"
	INVALID_TID_MSG               = "Terminal ID not found"
	INVALID_TOKEN_CODE            = "05"
	INVALID_TOKEN_MSG             = "Invalid access token"
	INVALID_SIGNATURE_CODE        = "30"
	INVALID_SIGNATURE_MSG         = "Invalid Signature"
	INVALID_TRX_CODE              = "12"
	INVALID_TRX_MSG               = "Invalid transaction"
	QR_EXPIRED_CODE               = "18"
	QR_EXPIRED_MSG                = "QR expired"
	BAD_REQUEST_CODE              = "30"
	BAD_REQUEST_MSG               = "Bad Request"
	INVALID_TIP_CODE              = "30"
	INVALID_TIP_MSG               = "Invalid tip indicator"
	BANK_NOT_SUPPORT_CODE         = "31"
	BANK_NOT_SUPPORT_MSG          = "Bank not supported by switch"
	INACTIVE_MERCHANT_CODE        = "31"
	INACTIVE_MERCHANT_MSG         = "Merchant inactive"
	BLOCKED_MERCHANT_CODE         = "32"
	BLOCKED_MERCHANT_MSG          = "Merchant blocked"
	MCC_NOT_MATCH_CODE            = "33"
	MCC_NOT_MATCH_MSG             = "MCC doesn't match"
	CRITERIA_NOT_MATCH_CODE       = "33"
	CRITERIA_NOT_MATCH_MSG        = "Merchant Criteria doesn't match"
	QR_NOT_FOUND_CODE             = "40"
	QR_NOT_FOUND_MSG              = "Transaction Not Found"
	INACTIVE_TID_CODE             = "41"
	INACTIVE_TID_MSG              = "Terminal ID inactive"
	AMOUNT_LIMIT_CODE             = "61"
	AMOUNT_LIMIT_MSG              = "Exceeds amount transactions limit"
	RESPONSE_LATE_CODE            = "68"
	RESPONSE_LATE_MSG             = "Response receive to late"
	GENERAL_ERROR_CODE            = "90"
	GENERAL_ERROR_MSG             = "General error"
	NOT_ALLOWED_CODE              = "98"
	NOT_ALLOWED_MSG               = "IP Address not authorized"
	INVALID_QR_TYPE_MSG           = "Invalid qr_type value"
	UNABLE_AUTHORIZE_CODE         = "99"
	UNABLE_AUTHORIZE_MSG          = "Unable to authorize"
	FAILED_CODE                   = "A0"
	FAILED_MSG                    = "Failed"
	TRANSACTION_NOT_FOUND_CODE    = "40"
	TRANSACTION_NOT_FOUND_MESSAGE = "Transaction Not Found"
)

func ErrRecordNotFound() CustomErrorResponse {
	return CustomErrorResponse{
		Message:  "Data tidak ditemukan.",
		ErrCode:  BAD_REQUEST_MSG,
		HTTPCode: http.StatusNotFound,
	}
}

func ErrUserNotFound() CustomErrorResponse {
	return CustomErrorResponse{
		Message:  "User tidak ditemukan.",
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
		HTTPCode: entity.StatusUnprocessableEntity,
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
		Message:  "Something went wrong.",
		ErrCode:  BAD_REQUEST_MSG,
		HTTPCode: entity.StatusUnprocessableEntity,
	}
}

func ErrInvalidRequest() CustomErrorResponse {
	return CustomErrorResponse{
		Message:  "Request data is invalid",
		ErrCode:  BAD_REQUEST_MSG,
		HTTPCode: entity.StatusUnprocessableEntity,
	}
}

func CustomError(message string, errCode string, httpCode int) CustomErrorResponse {
	return CustomErrorResponse{
		Message:  message,
		ErrCode:  errCode,
		HTTPCode: httpCode,
	}
}
