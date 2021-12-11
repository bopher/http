package session

import (
	"time"

	"github.com/bopher/cache"
	"github.com/gofiber/fiber/v2"
)

// NewCookieSession create new cookie based session
func NewCookieSession(
	cache cache.Cache,
	ctx *fiber.Ctx,
	secure bool,
	domain string,
	sameSite string,
	exp time.Duration,
	generator func() string,
	name string,
) Session {
	s := new(cSession)
	s.init(cache, ctx, secure, domain, sameSite, exp, generator, name)
	return s
}

// NewHeaderSession create new header based session
func NewHeaderSession(
	cache cache.Cache,
	ctx *fiber.Ctx,
	exp time.Duration,
	generator func() string,
	key string,
) Session {
	s := new(hSession)
	s.init(cache, ctx, exp, generator, key)
	return s
}
