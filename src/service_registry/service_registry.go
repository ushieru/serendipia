package serviceregistry

import (
	"fmt"
	"time"
)

type ServiceRegistry struct {
	services  map[string]service
	heartBeat int64
	handlers  []func()
}

func (serviceRegistry *ServiceRegistry) SetHandler(handler func()) {
	serviceRegistry.handlers = append(serviceRegistry.handlers, handler)
}

func (serviceRegistry *ServiceRegistry) Register(name string, ip string, port string) string {
	key := name + ip + port
	now := time.Now().UTC().UnixNano() / int64(time.Millisecond)
	if val, ok := serviceRegistry.services[key]; ok {
		val.timestamp = now
		fmt.Printf("[ServiceRegistry] Updated service: %s at %s:%s", name, ip, port)
		return key
	}
	service := service{
		name:      name,
		ip:        ip,
		port:      port,
		timestamp: now,
	}
	serviceRegistry.services[key] = service
	fmt.Printf("[ServiceRegistry] Registered service: %s at %s:%s", name, ip, port)
	return key
}

func (serviceRegistry *ServiceRegistry) Unregister(name string, ip string, port string) {
	key := name + ip + port
	delete(serviceRegistry.services, key)
	print("[ServiceRegistry] Unregister service: %s at %s:%s", name, ip, port)
}

func (serviceRegistry *ServiceRegistry) Cleanup() {
	now := time.Now().UTC().UnixNano() / int64(time.Millisecond)
	for key, service := range serviceRegistry.services {
		willBeEliminated := now-service.timestamp > int64(serviceRegistry.heartBeat)
		if willBeEliminated {
			delete(serviceRegistry.services, key)
			print("[ServiceRegistry] Eliminating service: %s", key)
		}
	}
}
