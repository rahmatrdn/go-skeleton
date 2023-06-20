package middleware

import (
	"github.com/gofiber/fiber/v2"
	apperr "gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/error"
	"gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/internal/http/auth"
)

func VerifyJWTToken(c *fiber.Ctx) error {
	if err := auth.VerifyToken(c); err != nil {
		return c.Status(apperr.ErrInvalidToken().HTTPCode).JSON(apperr.ErrInvalidToken())
	}

	return c.Next()
}
