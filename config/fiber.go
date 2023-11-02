package config

import (
	"github.com/gofiber/fiber/v2"
)

func NewFiberConfiguration() fiber.Config {
	return fiber.Config{
		CaseSensitive: true,
		ColorScheme: fiber.Colors{
			Black: "\u001b[39m",
		},
		StrictRouting: true,
		AppName:       "GO SKELETON v1.0.1",
		// ErrorHandler:  apperr.ErrorHandler,
	}
}
	