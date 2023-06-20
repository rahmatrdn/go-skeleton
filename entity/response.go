package entity

type GeneralResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

const (
	StatusUnprocessableEntity = 422
)

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
