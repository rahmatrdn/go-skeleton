package entity

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
	DATA_NOT_FOUND_MSG   = "Data not found"
	USER_NOT_FOUND_MSG   = "User not found"

	GENERAL_ERROR_MESSAGE = "Something went wrong. Please try again later."
)

type GeneralResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type CustomErrorResponse struct {
	Message  string `json:"message,omitempty"`
	ErrCode  string `json:"code,omitempty"`
	HTTPCode int    `json:"http_code"`
}
type CustomErrorResponseWithMeta struct {
	Message  string          `json:"message,omitempty"`
	ErrCode  string          `json:"code,omitempty"`
	HTTPCode int             `json:"http_code"`
	Meta     []ErrorResponse `json:"meta,omitempty"`
}
