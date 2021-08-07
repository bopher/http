package session

import "github.com/gofiber/fiber/v2"

// Session interface
type Session interface {
	// Parse session from request
	Parse()
	// ID get session id
	ID() string
	// Context get request context
	Context() *fiber.Ctx
	// Regenerate session id
	Regenerate()
	// Set session value
	Set(key string, value interface{})
	// Exists check if session is exists
	Exists(key string) bool
	// Get session value
	Get(key string) interface{}
	// Bool parse item as boolean
	Bool(key string, fallback bool) (bool, bool)
	// Int parse item as int
	Int(key string, fallback int) (int, bool)
	// Int8 parse item as int8
	Int8(key string, fallback int8) (int8, bool)
	// Int16 parse item as int16
	Int16(key string, fallback int16) (int16, bool)
	// Int32 parse item as int32
	Int32(key string, fallback int32) (int32, bool)
	// Int64 parse item as int64
	Int64(key string, fallback int64) (int64, bool)
	// UInt parse item as uint
	UInt(key string, fallback uint) (uint, bool)
	// UInt8 parse item as uint8
	UInt8(key string, fallback uint8) (uint8, bool)
	// UInt16 parse item as uint16
	UInt16(key string, fallback uint16) (uint16, bool)
	// UInt32 parse item as uint32
	UInt32(key string, fallback uint32) (uint32, bool)
	// UInt64 parse item as uint64
	UInt64(key string, fallback uint64) (uint64, bool)
	// Float32 parse item as float64
	Float32(key string, fallback float32) (float32, bool)
	// Float64 parse item as float64
	Float64(key string, fallback float64) (float64, bool)
	// String parse item as string
	String(key string, fallback string) (string, bool)
	// Bytes parse item as bytes array
	Bytes(key string, fallback []byte) ([]byte, bool)
	// All get all session stored value
	All() map[string]interface{}
	// Delete session value
	Delete(key string)
	// Destroy session
	Destroy()
	// Save session
	Save()
}
