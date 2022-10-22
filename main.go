package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Use("/*", func(c *fiber.Ctx) error {
		return c.SendString(c.Params("*"))
	})

	app.Listen(":3000")
}
