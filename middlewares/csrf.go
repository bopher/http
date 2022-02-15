package middlewares

import (
	"github.com/bopher/http/session"
	"github.com/bopher/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CSRFMiddleware protection middleware
func CSRFMiddleware(session session.Session) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if session == nil {
			return utils.TaggedError([]string{"CSRFMiddleware"}, "session driver is nil")
		}

		if !session.Exists("csrf_token") {
			session.Set("csrf_token", uuid.New().String())
		}

		return ctx.Next()
	}
}
