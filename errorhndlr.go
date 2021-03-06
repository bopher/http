package http

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/bopher/logger"
	"github.com/bopher/utils"
	"github.com/gofiber/fiber/v2"
)

// ErrorCallback a callback type for generate error response
type ErrorCallback func(*fiber.Ctx, error) error

// ErrorLogger handle errors and log into logger
//
// Enter only codes to log only codes included
func ErrorLogger(logger logger.Logger, formatter logger.TimeFormatter, callback ErrorCallback, onlyCodes ...int) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := 500
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		// Log
		if logger != nil && (len(onlyCodes) == 0 || utils.Contains[int](onlyCodes, code)) {
			logger.Divider("=", 100, c.IP())
			logger.Error().Tags(fmt.Sprintf("%d", code)).Print(err.Error())
			logger.Raw("\n")
			logger.Divider("-", 100, "Stacktrace:")
			logger.Raw(string(debug.Stack()))
			logger.Raw("\n")
			logger.Divider("-", 100, "Request Header:")
			logger.Raw(c.Request().Header.String())
			logger.Raw("\n")
			logger.Divider("-", 100, "Request Body:")
			logger.Raw(string(c.Request().Body()))
			logger.Raw("\n")
			logger.Divider("=", 100, formatter(time.Now().UTC(), "2006-01-02 15:04:05"))
			logger.Raw("\n\n")
		}

		// Return response
		if callback == nil {
			return c.SendStatus(code)
		} else {
			return callback(c, err)
		}
	}
}
