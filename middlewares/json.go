package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// JSONOnly implement
func JSONOnly(ctx *fiber.Ctx) error {
	if strings.ToLower(ctx.Get("Content-Type")) != "application/json" {
		return ctx.SendStatus(406)
	} else {
		return ctx.Next()
	}
}
