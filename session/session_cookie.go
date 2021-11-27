package session

import (
	"encoding/json"
	"time"

	"github.com/bopher/cache"
	"github.com/gofiber/fiber/v2"
)

type cookieSession struct {
	// cache driver
	cache cache.Cache
	// ctx request context
	ctx *fiber.Ctx
	// secure attribute for cookie
	secure bool
	// domain attribute for cookie
	domain string
	// sameSite attribute for cookie
	// Possible values: "Lax", "Strict", "None"
	// Optional. Default: "Lax"
	sameSite string
	// < 0 means when browser closes
	// > 0 is the time.Duration which the session cookies should expire.
	expiration time.Duration
	// Session id generator
	generator func() string

	key  string
	id   string
	data map[string]interface{}
}

func (s *cookieSession) init(cache cache.Cache, ctx *fiber.Ctx, secure bool, domain string, sameSite string, exp time.Duration, generator func() string, key string) {
	s.cache = cache
	s.ctx = ctx
	s.secure = secure
	s.domain = domain
	s.sameSite = sameSite
	s.expiration = exp
	s.generator = generator
	s.key = key
	if s.key == "" {
		s.key = "session"
	}
	s.data = make(map[string]interface{})
}

// identifier get cache identifier
func (s *cookieSession) identifier(id string) string {
	if id == "" {
		id = s.id
	}
	return "C-S-" + id
}

// Parse session from request
func (s *cookieSession) Parse() {
	sessionID := s.ctx.Cookies(s.key)
	if !s.cache.Exists(s.identifier(sessionID)) {
		s.Regenerate()
	} else {
		s.id = sessionID
		var res map[string]interface{}
		if data := s.cache.Bytes(s.identifier(""), nil); data != nil {
			if err := json.Unmarshal(data, &res); err == nil && res != nil {
				s.data = res
			}
		}
	}
}

// ID get session id
func (s *cookieSession) ID() string {
	return s.id
}

// Context get request context
func (s *cookieSession) Context() *fiber.Ctx {
	return s.ctx
}

// Regenerate session id
func (s *cookieSession) Regenerate() {
	// Clear old data
	s.Destroy()

	// Generate id
	s.id = s.generator()

	// Send Cookie
	cookie := fiber.Cookie{}
	cookie.Name = s.key
	cookie.Value = s.id
	cookie.Secure = s.secure
	cookie.Domain = s.domain
	cookie.SameSite = s.sameSite
	if s.expiration > 0 {
		cookie.Expires = time.Now().Add(s.expiration)
	}
	s.ctx.Cookie(&cookie)
}

// Set session value
func (s *cookieSession) Set(key string, value interface{}) {
	s.data[key] = value
}

// Exists check if session is exists
func (s *cookieSession) Exists(key string) bool {
	if _, ok := s.data[key]; ok {
		return true
	}
	return false
}

// Get session value
func (s *cookieSession) Get(key string) interface{} {
	if v, ok := s.data[key]; ok {
		return v
	}
	return nil
}

// Bool parse dependency as boolean
func (s *cookieSession) Bool(key string, fallback bool) (bool, bool) {
	if v, ok := s.data[key].(bool); ok {
		return v, true
	}
	return fallback, false
}

// Int parse dependency as int
func (s *cookieSession) Int(key string, fallback int) (int, bool) {
	if v, ok := s.data[key].(int); ok {
		return v, true
	}
	return fallback, false
}

// Int8 parse dependency as int8
func (s *cookieSession) Int8(key string, fallback int8) (int8, bool) {
	if v, ok := s.data[key].(int8); ok {
		return v, true
	}
	return fallback, false
}

// Int16 parse dependency as int16
func (s *cookieSession) Int16(key string, fallback int16) (int16, bool) {
	if v, ok := s.data[key].(int16); ok {
		return v, true
	}
	return fallback, false
}

// Int32 parse dependency as int32
func (s *cookieSession) Int32(key string, fallback int32) (int32, bool) {
	if v, ok := s.data[key].(int32); ok {
		return v, true
	}
	return fallback, false
}

// Int64 parse dependency as int64
func (s *cookieSession) Int64(key string, fallback int64) (int64, bool) {
	if v, ok := s.data[key].(int64); ok {
		return v, true
	}
	return fallback, false
}

// UInt parse dependency as uint
func (s *cookieSession) UInt(key string, fallback uint) (uint, bool) {
	if v, ok := s.data[key].(uint); ok {
		return v, true
	}
	return fallback, false
}

// UInt8 parse dependency as uint8
func (s *cookieSession) UInt8(key string, fallback uint8) (uint8, bool) {
	if v, ok := s.data[key].(uint8); ok {
		return v, true
	}
	return fallback, false
}

// UInt16 parse dependency as uint16
func (s *cookieSession) UInt16(key string, fallback uint16) (uint16, bool) {
	if v, ok := s.data[key].(uint16); ok {
		return v, true
	}
	return fallback, false
}

// UInt32 parse dependency as uint32
func (s *cookieSession) UInt32(key string, fallback uint32) (uint32, bool) {
	if v, ok := s.data[key].(uint32); ok {
		return v, true
	}
	return fallback, false
}

// UInt64 parse dependency as uint64
func (s *cookieSession) UInt64(key string, fallback uint64) (uint64, bool) {
	if v, ok := s.data[key].(uint64); ok {
		return v, true
	}
	return fallback, false
}

// Float32 parse dependency as float64
func (s *cookieSession) Float32(key string, fallback float32) (float32, bool) {
	if v, ok := s.data[key].(float32); ok {
		return v, true
	}
	return fallback, false
}

// Float64 parse dependency as float64
func (s *cookieSession) Float64(key string, fallback float64) (float64, bool) {
	if v, ok := s.data[key].(float64); ok {
		return v, true
	}
	return fallback, false
}

// String parse dependency as string
func (s *cookieSession) String(key string, fallback string) (string, bool) {
	if v, ok := s.data[key].(string); ok {
		return v, true
	}
	return fallback, false
}

// Bytes parse dependency as bytes array
func (s *cookieSession) Bytes(key string, fallback []byte) ([]byte, bool) {
	if v, ok := s.data[key].([]byte); ok {
		return v, true
	}
	return fallback, false
}

// All get all session stored value
func (s *cookieSession) All() map[string]interface{} {
	return s.data
}

// Delete session value
func (s *cookieSession) Delete(key string) {
	delete(s.data, key)
}

// Destroy session
func (s *cookieSession) Destroy() {
	s.cache.Forget(s.identifier(""))
	s.id = ""
	s.data = make(map[string]interface{})
}

// Save session
func (s *cookieSession) Save() {
	if data, err := json.Marshal(s.data); err == nil {
		if s.cache.Exists(s.identifier("")) {
			s.cache.Set(s.identifier(""), string(data))
		} else {
			if s.expiration <= 0 {
				s.cache.Put(s.identifier(""), string(data), 24*time.Hour)
			} else {
				s.cache.Put(s.identifier(""), string(data), s.expiration)
			}
		}
	}
}
