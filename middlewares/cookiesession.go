package middlewares

import (
	"time"

	"github.com/bopher/cache"
	"github.com/bopher/http/session"
	"github.com/gofiber/fiber/v2"
)

// NewCookieSession create new cookie based session
func NewCookieSession(cache cache.Cache, secure bool, domain string, sameSite string, exp time.Duration) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		s := session.NewCookieSession(cache, ctx, secure, domain, sameSite, exp, session.UUIDGenerator, "session")
		defer s.Save()
		s.Parse()
		ctx.Locals("c-session", s)
		return ctx.Next()
	}
}

// GetCookieSession get session driver from context
func GetCookieSession(ctx *fiber.Ctx) session.Session {
	if session, ok := ctx.Locals("c-session").(session.Session); ok {
		return session
	}
	return nil
}
