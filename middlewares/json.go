package middlewares

import (
	"github.com/bopher/http"
	"github.com/gofiber/fiber/v2"
)

// JSONOnly implement
func JSONOnly(ctx *fiber.Ctx) error {
	if !http.IsJsonRequest(ctx) {
		return ctx.SendStatus(fiber.StatusNotAcceptable)
	} else {
		return ctx.Next()
	}
}
