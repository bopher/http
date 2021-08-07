package middlewares

import (
	"github.com/bopher/cache"
	"github.com/gofiber/fiber/v2"
)

// Maintenance middleware
func Maintenance(c cache.Cache) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if c.Exists("maintenance") {
			return ctx.SendStatus(503)
		} else {
			return ctx.Next()
		}
	}
}
