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
NewCookieSession(cache cache.Cache, ctx *fiber.Ctx, secure bool, domain string, sameSite string, exp time.Duration, generator func() string, name string) Session

// Example:
import "github.com/bopher/http/session"
cSession := session.NewCookieSession(rCache, ctx, false, "", session.SameSiteLax, 30 * time.Minute, session.UUIDGenerator, "session")
```

### Create Header Session

Header sessions attached to and parsed from HTTP headers.

**Note:** If expiration time set to zero cache deleted after 24 hour.

```go
// Signature:
NewHeaderSession(cache cache.Cache, ctx *fiber.Ctx, exp time.Duration, generator func() string, name string) Session

// Example:
import "github.com/bopher/http/session"
hSession := NewHeaderSession(rCache, ctx, 30 * time.Minute, session.UUIDGenerator, "X-SESSION-ID")
```

### Usage

Session interface contains following methods:

#### ID

Get session id.

```go
ID() string
```

#### Context

Get request context.

```go
Context() *fiber.Ctx
```

#### Parse

Parse session from request.

```go
Parse() error
```

#### Regenerate

Regenerate session id.

```go
Regenerate() error
```

#### Set

Set session value.

```go
Set(key string, value interface{})
```

#### Get

Get session value.

```go
Get(key string) interface{}
```

#### Delete

Delete session value.

```go
Delete(key string)
```

#### Exists

Check if session is exists.

```go
Exists(key string) bool
```

#### Cast

Parse session item as caster.

```go
Cast(key string) caster.Caster
```

#### Destroy

Destroy session.

```go
Destroy() error
```

#### Save

Save session (must called at end of request).

```go
Save() error
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

**Note:** You can use `GetSession(ctx)` method for resolve session from cookie or session (if cookie not exists then try parse from header).

### CSRF Token

This middleware automatically generate and attach CSRF key to session.

```go
// Signature:
CSRFMiddleware(session session.Session) fiber.Handler

// Example:
import "github.com/bopher/http/middlewares"
app.Use(middlewares.CSRFMiddleware(mySession))

// Access CSRF key
if csrfKey, err := middlewares.GetCSRFKey(mySession); csrfKey != "" {
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

This middleware limit maximum request to server. this middleware send `X-LIMIT-UNTIL` header on locked and `X-LIMIT-REMAIN` otherwise.

```go
// Signature:
RateLimiter(key string, maxAttempts uint32, ttl time.Duration, c cache.Cache) fiber.Handler

// Example:
import "github.com/bopher/http/middlewares"
app.Use(middlewares.RateLimiter("global", 60, 1 * time.Minute, rCache)) // Accept 60 request in minutes
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
