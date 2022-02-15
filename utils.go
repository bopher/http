package http

import (
	"strings"

	"github.com/bopher/cache"
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
	return ctx.Accepts("application/json") == "application/json"
}

// IsUnderMaintenance check if under maintenance mode
func IsUnderMaintenance(c cache.Cache) (bool, error) {
	return c.Exists("maintenance")
}

// GetCSRF get csrf key
func GetCSRF(session session.Session) (string, error) {
	if session == nil {
		return "", utils.TaggedError([]string{"GetCSRF"}, "session driver is nil")
	}

	caster := session.Cast("csrf_token")
	if caster.IsNil() {
		return "", nil
	}

	v, err := caster.String()
	if err != nil {
		return "", utils.TaggedError([]string{"GetCSRF"}, err.Error())
	}

	return v, nil
}

// CheckCSRF check csrf token
func CheckCSRF(session session.Session, key string) (bool, error) {
	if k, err := GetCSRF(session); err != nil {
		return false, err
	} else {
		return k != "" && k == key, nil
	}
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
