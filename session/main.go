package session

import (
	"time"

	"github.com/bopher/cache"
	"github.com/gofiber/fiber/v2"
)

// NewCookieSession create new cookie based session
func NewCookieSession(cache cache.Cache, ctx *fiber.Ctx, secure bool, domain string, sameSite string, exp time.Duration, generator func() string, key string) Session {
	s := new(cookieSession)
	s.init(cache, ctx, secure, domain, sameSite, exp, generator, key)
	return s
}

// NewHeaderSession create new header based session
func NewHeaderSession(cache cache.Cache, ctx *fiber.Ctx, exp time.Duration, generator func() string, key string) Session {
	s := new(headerSession)
	s.init(cache, ctx, exp, generator, key)
	return s
}
