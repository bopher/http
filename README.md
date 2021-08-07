# HTTP

Session manager, middlewares and error handler for gofiber app.

## Session

Http packages comes with two type session driver by default (header and cookie).

### Requirements Knowledge

- Session use `"github.com/bopher/cache"` for storing session data.
- Session required a generator function `func() string` for creating unique session id. By default session driver contains UUID generator function.

### Create Cookie Session

**Note:** Set expiration 0 to ignore cookie expiration time (delete cookie on browser close and delete cache data after 24 hour).

```go
// Signature:
NewCookieSession(cache cache.Cache, ctx *fiber.Ctx, secure bool, domain string, sameSite string, exp time.Duration, generator func() string, key string) Session

// Example:
import "github.com/bopher/http/session"
cSession := session.NewCookieSession(rCache, ctx, false, "", session.SameSiteLax, 30 * time.Minute, session.UUIDGenerator, "session")
```

### Create Header Session

Header sessions attached to and parsed from HTTP headers.

**Note:** If expiration time set to zero cache deleted after 24 hour.

```go
// Signature:
NewHeaderSession(cache cache.Cache, ctx *fiber.Ctx, exp time.Duration, generator func() string, key string) Session

// Example:
import "github.com/bopher/http/session"
hSession := NewHeaderSession(rCache, ctx, 30 * time.Minute, session.UUIDGenerator, "X-SESSION-ID")
```

### Usage

Session interface contains following methods:

#### Parse

Parse session from request.

```go
// Signature:
Parse()
```

#### ID

Get session id.

```go
// Signature:
ID() string
```

#### Context

Get request context.

```go
// Signature:
Context() *fiber.Ctx
```

#### Regenerate

Regenerate session id.

```go
// Signature:
Regenerate()
```

#### Set

Set session item.

```go
// Signature:
Set(key string, value interface{})
```

#### Exists

Check if session item is exists.

```go
// Signature:
Exists(key string) bool
```

#### Get

Get session item.

```go
// Signature:
Get(key string) interface{}
```

#### All

Get all session stored value.

```go
// Signature:
All() map[string]interface{}
```

#### Delete

Delete session item.

```go
// Signature:
Delete(key string)
```

#### Destroy

Destroy session.

```go
// Signature:
Destroy()
```

#### Save

Save session (must called at ned of request).

```go
// Signature:
Save()
```

#### Get By Type Methods

Helper get methods return fallback value if value not exists in session.

```go
// Parse item as boolean
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
```

## Middlewares

HTTP Package contains following middlewares by default:

### Cookie Session

This middleware create a session from cookie.

```go
// Signature:
NewCookieSession(cache cache.Cache, secure bool, domain string, sameSite string, exp time.Duration) fiber.Handler

// Example:
import "github.com/bopher/http/middlewares"
import "github.com/bopher/http/session"
app.Use(middlewares.NewCookieSession(myCache, false, "", session.SameSiteNone, 0))

// Access session
session := middlewares.GetCookieSession(ctx)
if session != nil {
    // Do something with session
}
```

### Header Session

This middleware create a session from HTTP header.

```go
// Signature:
NewHeaderSession(cache cache.Cache, exp time.Duration) fiber.Handler

// Example:
import "github.com/bopher/http/middlewares"
import "github.com/bopher/http/session"
app.Use(middlewares.NewHeaderSession(myCache, 0))

// Access session
session := middlewares.GetHeaderSession(ctx)
if session != nil {
    // Do something with session
}
```

### CSRF Token

This middleware automatically generate and attach CSRF key to session.

```go
// Signature:
CSRFMiddleware(session session.Session) fiber.Handler

// Example:
import "github.com/bopher/http/middlewares"
app.Use(middlewares.CSRFMiddleware(mySession))

// Access CSRF key
if csrfKey := middlewares.GetCSRFKey(mySession); csrfKey != "" {
    // check CSRF key
}
```

### JSON Only Checker

This middleware returns 406 HTTP error if request not want json (check `Content-Type` header). This middlewares useful for api requests.

```go
// Signature:
JSONOnly(ctx *fiber.Ctx) error

// Example:
import "github.com/bopher/http/middlewares"
app.Use(middlewares.JSONOnly)
```

### Maintenance Mode

This middleware return 503 HTTP error if `maintenance` key exists in cache.

```go
// Signature:
Maintenance(c cache.Cache) fiber.Handler

// Example:
import "github.com/bopher/http/middlewares"
app.Use(middlewares.Maintenance(rCache))
```

### Access Logger

This middleware format and log request information to logger (`"github.com/bopher/logger"` driver).

```go
// Signature:
AccessLogger(logger logger.Logger) fiber.Handler

// Example:
import "github.com/bopher/http/middlewares"
app.Use(middlewares.AccessLogger(myLogger))
```

### Rate Limiter

This middleware limit maximum request to server.

```go
// Signature:
RateLimiter(key string, maxAttempts uint32, ttl time.Duration, c cache.Cache) fiber.Handler

// Example:
import "github.com/bopher/http/middlewares"
app.Use(middlewares.RateLimiter("req-rl", 60, 1 * time.Minute, rCache)) // Accept 60 request in minutes
```

## Recover Panics (Fiber ErrorHandler)

This Error handler log error to logger and return http error to response. You can use this function instead of default fiber error handler.

**Note:** You can pass a list of code as _onlyCodes_ parameter to log errors only if error code contains in your list.

```go
// Signature:
ErrorLogger(logger logger.Logger, formatter logger.TimeFormatter, onlyCodes ...int) fiber.ErrorHandler

// Example:
import "github.com/bopher/http"
app := fiber.New(fiber.Config{
    ErrorHandler: http.ErrorLogger(myLogger, myFormatter),
})
```
