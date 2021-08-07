package middlewares

import (
	"errors"

	"github.com/bopher/http/session"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CSRFMiddleware protection middleware
func CSRFMiddleware(session session.Session) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if session == nil {
			return errors.New("session driver not valid!")
		}

		token, _ := session.Get("csrf_token").(string)
		if token == "" {
			session.Set("csrf_token", uuid.New().String())
		}

		return ctx.Next()
	}
}

// CSRFMiddleware
func GetCSRFKey(session session.Session) string {
	if session != nil {
		if token, ok := session.Get("csrf_token").(string); ok {
			return token
		}
	}
	return ""
}
