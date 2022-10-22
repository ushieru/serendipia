package serviceregistry

import (
	"errors"
	"fmt"
	"github.com/ushieru/serendipia/src/utils"
	"time"
)

type ServiceRegistry struct {
	Services  map[string]Service
	HeartBeat int64
	Handlers  []func()
}

func NewServiceRegistry(heartBeat int64) *ServiceRegistry {
	return &ServiceRegistry{
		Services:  make(map[string]Service),
		HeartBeat: heartBeat,
		Handlers:  make([]func(), 0),
	}
}

func (serviceRegistry *ServiceRegistry) SetHandler(handler func()) {
	serviceRegistry.Handlers = append(serviceRegistry.Handlers, handler)
}

func (serviceRegistry *ServiceRegistry) Register(name string, ip string, port string) string {
	key := name + ip + port
	now := utils.GetTimeStamp()
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

func (serviceRegistry *ServiceRegistry) Get(name string) (*Service, error) {
	for _, service := range serviceRegistry.Services {
		if service.Name == name {
			return &service, nil
		}
	}
	return nil, errors.New("Service not found")
}

func (serviceRegistry *ServiceRegistry) Unregister(name string, ip string, port string) {
	key := name + ip + port
	delete(serviceRegistry.Services, key)
	print("[ServiceRegistry] Unregister service: %s at %s:%s", name, ip, port)
}

func (serviceRegistry *ServiceRegistry) Cleanup() {
	now := utils.GetTimeStamp()
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
