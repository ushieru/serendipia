package serviceregistry

import (
	"fmt"
	"time"
)

type ServiceRegistry struct {
	Services  map[string]Service
	HeartBeat int64
	Handlers  []func()
}

func (serviceRegistry *ServiceRegistry) SetHandler(handler func()) {
	serviceRegistry.Handlers = append(serviceRegistry.Handlers, handler)
}

func (serviceRegistry *ServiceRegistry) Register(name string, ip string, port string) string {
	key := name + ip + port
	now := time.Now().UTC().UnixNano() / int64(time.Millisecond)
	if val, ok := serviceRegistry.Services[key]; ok {
		val.Timestamp = now
		fmt.Printf("[ServiceRegistry] Updated service: %s at %s:%s", name, ip, port)
		return key
	}
	service := Service{
		Name:      name,
		Ip:        ip,
		Port:      port,
		Timestamp: now,
	}
	serviceRegistry.Services[key] = service
	fmt.Printf("[ServiceRegistry] Registered service: %s at %s:%s", name, ip, port)
	return key
}

func (serviceRegistry *ServiceRegistry) Unregister(name string, ip string, port string) {
	key := name + ip + port
	delete(serviceRegistry.Services, key)
	print("[ServiceRegistry] Unregister service: %s at %s:%s", name, ip, port)
}

func (serviceRegistry *ServiceRegistry) Cleanup() {
	now := time.Now().UTC().UnixNano() / int64(time.Millisecond)
	for key, service := range serviceRegistry.Services {
		willBeEliminated := now-service.Timestamp > int64(serviceRegistry.HeartBeat)
		if willBeEliminated {
			delete(serviceRegistry.Services, key)
			print("[ServiceRegistry] Eliminating service: %s", key)
		}
	}
}

func (serviceRegistry *ServiceRegistry) Init() {
	for {
		serviceRegistry.Cleanup()
		time.Sleep(time.Duration(serviceRegistry.HeartBeat) * time.Second)
	}
}
