package middlewares

import (
	"time"

	"github.com/bopher/cache"
	"github.com/bopher/http/session"
	"github.com/gofiber/fiber/v2"
)

// NewHeaderSession create new header based session
func NewHeaderSession(cache cache.Cache, exp time.Duration) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		s := session.NewHeaderSession(cache, ctx, exp, session.UUIDGenerator, "X-SESSION-ID")
		defer s.Save()
		s.Parse()
		ctx.Locals("h-session", s)
		return ctx.Next()
	}
}

// GetHeaderSession get session driver from context
func GetHeaderSession(ctx *fiber.Ctx) session.Session {
	if session, ok := ctx.Locals("h-session").(session.Session); ok {
		return session
	}
	return nil
}
