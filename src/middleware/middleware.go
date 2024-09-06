package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Outputs a log message for each request recieved in the format:
//
// "[TIME] [DATE] INC REQ: [METHOD] [PATH] at [TIME] from [IP] response: [STATUS_CODE]"
func LogRequests(c *fiber.Ctx) error {
	log.Printf("[REQ LOG]: %s %-20s at %s origin %s result %d\n",
		c.Method(),
		c.Path(),
		time.Now().Format(time.RFC3339),
		c.IP(),
		c.Response().StatusCode())
	return nil
}
