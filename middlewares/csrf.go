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
			return utils.TaggedError(
				[]string{"CSRFMiddleware"},
				"session driver is nil",
			)
		}

		if session.Get("csrf_token") == nil {
			session.Set("csrf_token", uuid.New().String())
		}

		return ctx.Next()
	}
}

// CSRFMiddleware
func GetCSRFKey(session session.Session) (string, error) {
	if session == nil {
		return "", utils.TaggedError(
			[]string{"GetCSRFKey"},
			"session driver is nil",
		)
	}

	caster := session.Cast("csrf_token")
	if caster.IsNil() {
		return "", nil
	}

	v, err := caster.String()
	if err != nil {
		return "", utils.TaggedError(
			[]string{"GetCSRFKey"},
			err.Error(),
		)
	}

	return v, nil
}
