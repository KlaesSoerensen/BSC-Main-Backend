package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// Outputs a log message for each request recieved in the format:
//
// "[TIME] [DATE] INC REQ: [Remote IP] --> [Method] [url] --> [STATUS_CODE]"
func LogRequests(c *fiber.Ctx) error {
	log.Printf("[REQ LOG] %s --> %-7s %-30s --> %d\n",
		c.IP(),
		c.Method(),
		c.Path(),
		c.Response().StatusCode())
	return nil
}
