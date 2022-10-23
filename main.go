package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ushieru/serendipia/src/gateway"
	"github.com/ushieru/serendipia/src/service_registry"
)

func main() {
	app := fiber.New()
	proxy := gateway.NewGateway()
	go proxy.ServiceRegistry.InitTick()

	app.Get("/services", func(ctx *fiber.Ctx) error {
		services := make(map[string]serviceregistry.Service)
		proxy.ServiceRegistry.Services.Range(func(key, value any) bool {
			services[key.(string)] = value.(serviceregistry.Service)
			return true
		})
		return ctx.JSON(services)
	})

	app.Post("/services", func(ctx *fiber.Ctx) error {
		params := &struct {
			ServiceName string `json:"service_name"`
			ServicePort string `json:"service_port"`
		}{}
		if err := ctx.BodyParser(&params); err != nil {
			return ctx.SendStatus(400)
		}
		serviceKey := proxy.
			ServiceRegistry.
			Register(params.ServiceName, ctx.IP(), ctx.Protocol(), params.ServicePort)
		return ctx.SendString(serviceKey)
	})

	app.Use("/*", func(c *fiber.Ctx) error {
		return proxy.CallService(c)
	})

	app.Listen(":3000")
}
