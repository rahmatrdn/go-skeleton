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

type Presenter struct{}

// NewPresenter initialize new JSON presenter that used to hold logic for presenter logic
func NewPresenter() *Presenter {
	return &Presenter{}
}

// SuccessBody is used to define success response body data structure
type ResponseBody struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Code    string      `json:"code"`
}

func (p *Presenter) BuildSuccess(c *fiber.Ctx, data interface{}, message string, code int) error {
	response := &ResponseBody{
		Data:    data,
		Message: message,
		Code:    apperr.SUCCESS_CODE,
	}

	return c.JSON(response)
}

func (p *Presenter) BuildError(c *fiber.Ctx, err error) error {
	unwrappedErr := errors.Unwrap(err)

	if unwrappedErr != nil {
		errorData := strings.Split(unwrappedErr.Error(), "XX: ")
		errorCode := errorData[1]
		errorMessage := errorData[0]

		// Handle error struct validation
		if errorCode == apperr.INVALID_PAYLOAD_CODE {
			var errResponse []entity.ErrorResponse
			json.Unmarshal([]byte(errorMessage), &errResponse)

			return c.Status(apperr.ErrGeneralInvalid().HTTPCode).
				JSON(
					apperr.ErrInvalidPayload(errResponse),
				)
		}
	}

	switch err.(type) {
	case apperr.CustomErrorResponse:
		httpCode := err.(apperr.CustomErrorResponse).HTTPCode

		return c.Status(httpCode).JSON(err)
	default:
		return c.Status(apperr.ErrGeneralInvalid().HTTPCode).
			JSON(apperr.CustomError(err.Error(),
				apperr.BAD_REQUEST_CODE,
				http.StatusUnprocessableEntity))
	}
}
