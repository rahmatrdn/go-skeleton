package error

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/entity"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {

	_, validationError := err.(ValidationError)
	if validationError {
		data := err.Error()
		var messages []map[string]interface{}

		errJson := json.Unmarshal([]byte(data), &messages)
		PanicLogging(errJson)
		return ctx.JSON(entity.GeneralResponse{
			Code:    400,
			Message: "Bad Request",
			Data:    messages,
		})
	}

	_, notFoundError := err.(NotFoundError)
	if notFoundError {
		return ctx.JSON(entity.GeneralResponse{
			Code:    404,
			Message: "Not Found",
			Data:    err.Error(),
		})
	}

	_, unauthorizedError := err.(UnauthorizedError)
	if unauthorizedError {
		return ctx.JSON(entity.GeneralResponse{
			Code:    401,
			Message: "Unauthorized",
			Data:    err.Error(),
		})
	}

	_, internalError := err.(InternalError)
	if internalError {
		return ctx.JSON(entity.GeneralResponse{
			Code:    500,
			Message: "General Error",
			Data:    err.Error(),
		})
	}

	return ctx.JSON(entity.GeneralResponse{
		Code:    500,
		Message: "General Error",
		Data:    err.Error(),
	})
}

func PanicLogging(err interface{}) {
	if err != nil {
		panic(err)
	}
}
