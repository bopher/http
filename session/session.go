package session

import (
	"github.com/bopher/caster"
	"github.com/gofiber/fiber/v2"
)

// Session interface
type Session interface {
	// ID get session id
	ID() string
	// Context get request context
	Context() *fiber.Ctx
	// Parse parse session from request
	Parse() error
	// Regenerate regenerate session id
	Regenerate() error
	// Set set session value
	Set(key string, value any)
	// Get get session value
	Get(key string) any
	// Delete delete session value
	Delete(key string)
	// Exists check if session is exists
	Exists(key string) bool
	// Cast parse session item as caster
	Cast(key string) caster.Caster
	// Destroy session
	Destroy() error
	// Save session
	Save() error
}
