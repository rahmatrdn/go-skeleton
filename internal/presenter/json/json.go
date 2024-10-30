package json

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/rahmatrdn/go-skeleton/entity"
	apperr "github.com/rahmatrdn/go-skeleton/error"

	"github.com/gofiber/fiber/v2"
)

type Json struct{}

// NewPresenter initialize new JSON presenter that used to hold logic for presenter logic
func NewJsonPresenter() *Json {
	return &Json{}
}

type JsonPresenter interface {
	BuildSuccess(c *fiber.Ctx, data interface{}, message string, code int) error
	BuildError(c *fiber.Ctx, err error) error
}

// SuccessBody is used to define success response body data structure
type ResponseBody struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Code    string      `json:"code"`
}

func (p *Json) BuildSuccess(c *fiber.Ctx, data interface{}, message string, code int) error {
	response := &ResponseBody{
		Data:    data,
		Message: message,
		Code:    entity.SUCCESS_CODE,
	}

	return c.JSON(response)
}

func (p *Json) BuildError(c *fiber.Ctx, err error) error {
	unwrappedErr := errors.Unwrap(err)

	if unwrappedErr != nil {
		errorData := strings.Split(unwrappedErr.Error(), "XX: ")

		if len(errorData) < 2 {
			return c.Status(apperr.ErrGeneralInvalid().HTTPCode).
				JSON(apperr.CustomError(err.Error(),
					entity.BAD_REQUEST_CODE,
					http.StatusUnprocessableEntity))
		}

		errorCode := errorData[1]
		errorMessage := errorData[0]

		// Handle error struct validation
		if errorCode == entity.INVALID_PAYLOAD_CODE {
			var errResponse []entity.ErrorResponse
			json.Unmarshal([]byte(errorMessage), &errResponse)

			return c.Status(apperr.ErrGeneralInvalid().HTTPCode).
				JSON(
					apperr.ErrInvalidPayload(errResponse),
				)
		}
	}

	switch err := err.(type) {
	case apperr.CustomErrorResponse:
		httpCode := err.HTTPCode
		return c.Status(httpCode).JSON(err)
	default:
		return c.Status(apperr.ErrGeneralInvalid().HTTPCode).
			JSON(apperr.CustomError(err.Error(),
				entity.BAD_REQUEST_CODE,
				http.StatusUnprocessableEntity))
	}
}
