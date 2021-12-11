package middlewares

import (
	"time"

	"github.com/bopher/cache"
	"github.com/bopher/http/session"
	"github.com/gofiber/fiber/v2"
)

// GetCookieSession get session driver from context
func GetCookieSession(ctx *fiber.Ctx) session.Session {
	if session, ok := ctx.Locals("sessioncookie").(session.Session); ok {
		return session
	}
	return nil
}

// GetHeaderSession get session driver from context
func GetHeaderSession(ctx *fiber.Ctx) session.Session {
	if session, ok := ctx.Locals("sessionheader").(session.Session); ok {
		return session
	}
	return nil
}

// GetSession get session driver from context
//
// if cookie session exists return cookie session otherwise try to resolve header session or return nil on fail
func GetSession(ctx *fiber.Ctx) session.Session {
	if session := GetCookieSession(ctx); session != nil {
		return session
	} else {
		return GetHeaderSession(ctx)
	}
}

// NewCookieSession create new cookie based session
//
// this function generate panic on save fail!
func NewCookieSession(
	cache cache.Cache,
	secure bool,
	domain string,
	sameSite string,
	exp time.Duration,
) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		cSession := session.NewCookieSession(
			cache,
			ctx,
			secure,
			domain,
			sameSite,
			exp,
			session.UUIDGenerator,
			"session",
		)

		defer func(ses session.Session) {
			if ses == nil {
				panic("[CookieSessionMW] session is null!")
			}
			err := ses.Save()
			if err != nil {
				panic(err.Error())
			}
		}(cSession)

		if err := cSession.Parse(); err != nil {
			return err
		}

		ctx.Locals("sessioncookie", cSession)
		return ctx.Next()
	}
}

// NewHeaderSession create new header based session
//
// this function generate panic on save fail!
func NewHeaderSession(
	cache cache.Cache,
	exp time.Duration,
) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		hSession := session.NewHeaderSession(
			cache,
			ctx,
			exp,
			session.UUIDGenerator,
			"X-SESSION-ID",
		)

		defer func(ses session.Session) {
			if ses == nil {
				panic("[HeaderSessionMW] session is null!")
			}
			err := ses.Save()
			if err != nil {
				panic(err.Error())
			}
		}(hSession)

		if err := hSession.Parse(); err != nil {
			return err
		}

		ctx.Locals("sessionheader", hSession)
		ctx.Append("Access-Control-Expose-Headers", "X-SESSION-ID")
		ctx.Append("Access-Control-Allow-Headers", "X-SESSION-ID")
		return ctx.Next()
	}
}
