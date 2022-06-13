package middlewares

import (
	"fmt"
	"regexp"
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
	methods []string,
	ignore []string,
) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		prettyErr := func(err error) error {
			return utils.TaggedError(
				[]string{"RateLimiterMW"},
				err.Error(),
			)
		}

		validMethod := func(_method string) bool {
			if len(methods) == 0 || utils.Contains[string](methods, _method) {
				return true
			}
			return false
		}

		mustIgnore := func(path string) bool {
			if len(ignore) > 0 {
				for _, expr := range ignore {
					if r, err := regexp.Compile(expr); err != nil {
						if r.MatchString(path) {
							return true
						}
					}
				}
			}
			return false
		}

		if !validMethod(ctx.Method()) || mustIgnore(ctx.Path()) {
			return ctx.Next()
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
