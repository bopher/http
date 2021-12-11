package middlewares

import (
	"github.com/bopher/cache"
	"github.com/bopher/utils"
	"github.com/gofiber/fiber/v2"
)

// Maintenance middleware
func Maintenance(c cache.Cache) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		exists, err := c.Exists("maintenance")
		if err != nil {
			return utils.TaggedError(
				[]string{"MaintenanceMW"},
				err.Error(),
			)
		}

		if exists {
			return ctx.SendStatus(fiber.StatusServiceUnavailable)
		}
		return ctx.Next()
	}
}
