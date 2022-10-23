package gateway

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/ushieru/serendipia/src/circuit_breaker"
	"github.com/ushieru/serendipia/src/service_registry"
	"io"
	"strings"
)

type Gateway struct {
	ServiceRegistry serviceregistry.ServiceRegistry
	CircuitBreaker  circuitbreaker.CircuitBreaker
}

func NewGateway() *Gateway {
	return &Gateway{
		ServiceRegistry: *serviceregistry.NewServiceRegistry(5),
		CircuitBreaker:  *circuitbreaker.NewCircuitBreaker(5, 10, 2),
	}
}

func (gateway Gateway) CallService(c *fiber.Ctx) error {
	method := c.Method()
	paths := strings.Split(c.Params("*"), "/")
	serviceName, path := paths[0], strings.Join(paths[1:], "/")
	service, err := gateway.ServiceRegistry.Get(serviceName)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	url := service.Protocol + "://" + service.Ip + ":" + service.Port + "/" + path
	serviceResponse, err := gateway.CircuitBreaker.CallService(method, url, c.Context().RequestBodyStream(), c.GetReqHeaders())
	if err != nil {
		log.Println("[Gateway] Error >>>", err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer serviceResponse.Body.Close()
	body, responseBodyErr := io.ReadAll(serviceResponse.Body)
	if responseBodyErr != nil {
		log.Println("[Gateway] Response Body Error >>>", responseBodyErr.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	response := c.Response()
	response.SetBody(body)
	for key := range serviceResponse.Header {
		response.Header.Add(key, serviceResponse.Header.Get(key))
	}
	response.SetStatusCode(serviceResponse.StatusCode)
	return nil
}
