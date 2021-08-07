package middlewares

import (
	"strconv"
	"time"

	"github.com/bopher/cache"
	"github.com/gofiber/fiber/v2"
)

// RateLimiter middleware
func RateLimiter(key string, maxAttempts uint32, ttl time.Duration, c cache.Cache) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		limiter := cache.NewRateLimiter(key+"-R-L-"+ctx.IP(), maxAttempts, ttl, c)
		if limiter.MustLock() {
			ctx.Set("X-LIMIT-UNTIL", limiter.AvailableIn().String())
			return ctx.SendStatus(429)
		} else {
			limiter.Hit()
			ctx.Set("X-LIMIT-REMAIN", strconv.Itoa(int(limiter.RetriesLeft())))
			return ctx.Next()
		}
	}
}
