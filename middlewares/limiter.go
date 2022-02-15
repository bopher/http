package middlewares

import (
	"fmt"
	"time"

	"github.com/bopher/cache"
	"github.com/bopher/utils"
	"github.com/gofiber/fiber/v2"
)

// RateLimiter middleware
func RateLimiter(
	key string,
	maxAttempts uint32,
	ttl time.Duration,
	c cache.Cache,
	callback fiber.Handler,
) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		prettyErr := func(err error) error {
			return utils.TaggedError(
				[]string{"RateLimiterMW"},
				err.Error(),
			)
		}

		limiter, err := cache.NewRateLimiter(key+"_limiter_-"+ctx.IP(), maxAttempts, ttl, c)
		if err != nil {
			return prettyErr(err)
		}

		ctx.Append("Access-Control-Expose-Headers", "X-LIMIT-UNTIL")
		ctx.Append("Access-Control-Expose-Headers", "X-LIMIT-REMAIN")
		ctx.Append("Access-Control-Allow-Headers", "X-LIMIT-UNTIL")
		ctx.Append("Access-Control-Allow-Headers", "X-LIMIT-REMAIN")

		mustLook, err := limiter.MustLock()
		if err != nil {
			return prettyErr(err)
		}

		if mustLook {
			until, err := limiter.AvailableIn()
			if err != nil {
				return prettyErr(err)
			}
			ctx.Set("X-LIMIT-UNTIL", until.String())
			if callback == nil {
				return ctx.SendStatus(fiber.StatusTooManyRequests)
			} else {
				return callback(ctx)
			}
		} else {
			err = limiter.Hit()
			if err != nil {
				return prettyErr(err)
			}

			left, err := limiter.RetriesLeft()
			if err != nil {
				return prettyErr(err)
			}
			ctx.Set("X-LIMIT-REMAIN", fmt.Sprint(left))

			return ctx.Next()
		}
	}
}
