package middlewares

import (
	"github.com/bopher/http"
	"github.com/gofiber/fiber/v2"
)

// JSONOnly allow json requests only
func JSONOnly(callback fiber.Handler) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if !http.IsJsonRequest(ctx) {
			if callback == nil {
				return ctx.SendStatus(fiber.StatusNotAcceptable)
			} else {
				return callback(ctx)
			}
		} else {
			return ctx.Next()
		}
	}
}
