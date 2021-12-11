package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// JSONOnly implement
func JSONOnly(ctx *fiber.Ctx) error {
	if strings.ToLower(ctx.Get("Content-Type")) != "application/json" {
		return ctx.SendStatus(fiber.StatusNotAcceptable)
	} else {
		return ctx.Next()
	}
}
