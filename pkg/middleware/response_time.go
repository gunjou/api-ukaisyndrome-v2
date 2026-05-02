package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func ResponseTime() fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Locals("start_time", time.Now())

		return c.Next()
	}
}

