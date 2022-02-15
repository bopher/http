package http

import (
	"strings"

	"github.com/bopher/http/session"
	"github.com/bopher/utils"
	"github.com/gofiber/fiber/v2"
)

// IsJsonRequest check if request is json
func IsJsonRequest(ctx *fiber.Ctx) bool {
	return strings.ToLower(ctx.Get("Content-Type")) == "application/json"
}

// WantJson check if request want json
func WantJson(ctx *fiber.Ctx) bool {
	return strings.Contains(strings.ToLower(ctx.Get("Accept")), "application/json")
}

// CookieSession get cookie session driver from context
func CookieSession(ctx *fiber.Ctx) session.Session {
	if session, ok := ctx.Locals("sessioncookie").(session.Session); ok {
		return session
	}
	return nil
}

// HeaderSession get header session driver from context
func HeaderSession(ctx *fiber.Ctx) session.Session {
	if session, ok := ctx.Locals("sessionheader").(session.Session); ok {
		return session
	}
	return nil
}

// GetSession get session driver from context
//
// if cookie session exists return cookie session otherwise try to resolve header session or return nil on fail
func GetSession(ctx *fiber.Ctx) session.Session {
	if session := CookieSession(ctx); session != nil {
		return session
	} else {
		return HeaderSession(ctx)
	}
}

// GetCSRF get csrf key
func GetCSRF(ctx *fiber.Ctx) (string, error) {
	if ses := GetSession(ctx); ses == nil {
		return "", utils.TaggedError([]string{"GetCSRF"}, "session driver is nil")
	} else {
		return ses.Cast("csrf_token").StringSafe(""), nil
	}
}

// CheckCSRF check csrf token
func CheckCSRF(ctx *fiber.Ctx, key string) (bool, error) {
	if k, err := GetCSRF(ctx); err != nil {
		return false, err
	} else {
		return k != "" && k == key, nil
	}
}
