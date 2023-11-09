package config

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func NewFiberConfiguration(cfg *Config) fiber.Config {
	return fiber.Config{
		CaseSensitive: true,
		ColorScheme: fiber.Colors{
			Black: "\u001b[39m",
		},
		StrictRouting: true,
		AppName:       fmt.Sprintf("%s - %s", cfg.AppName, cfg.AppVersion),
	}
}
