package serviceregistry

import (
	"errors"
	"fmt"
	"github.com/ushieru/serendipia/src/utils"
	"sync"
	"time"
)

type ServiceRegistry struct {
	Services  sync.Map
	HeartBeat int64
}

func NewServiceRegistry(heartBeat int64) *ServiceRegistry {
	return &ServiceRegistry{
		HeartBeat: heartBeat,
	}
}

func (serviceRegistry *ServiceRegistry) Register(name string, ip string, protocol string, port string) string {
	key := name + ip + port
	now := utils.GetTimeStamp()
	if service, ok := serviceRegistry.Services.Load(key); ok {
		s := service.(Service)
		serviceRegistry.Services.Store(key, Service{
			Name:      s.Name,
			Ip:        s.Ip,
			Protocol:  s.Protocol,
			Port:      s.Port,
			Timestamp: now,
		})
		fmt.Println("[ServiceRegistry] Updated service:", name, "at", ip, ":", port, "\ntime", now)
		return key
	}
	service := Service{
		Name:      name,
		Protocol:  protocol,
		Ip:        ip,
		Port:      port,
		Timestamp: now,
	}
	serviceRegistry.Services.Store(key, service)
	fmt.Println("[ServiceRegistry] Registered service:", name, "at", ip, ":", port)
	return key
}

func (serviceRegistry *ServiceRegistry) Get(name string) (*Service, error) {
	var serviceResponse Service
	serviceRegistry.Services.Range(func(k, service interface{}) bool {
		if service.(Service).Name == name {
			serviceResponse = service.(Service)
			return false
		}
		return true
	})
	return &serviceResponse, errors.New("Service not found")
}

func (serviceRegistry *ServiceRegistry) Unregister(name string, ip string, port string) {
	key := name + ip + port
	serviceRegistry.Services.Delete(key)
	fmt.Println("[ServiceRegistry] Unregister service:", name, "at", ip, ":", port)
}

func (serviceRegistry *ServiceRegistry) Cleanup() {
	now := utils.GetTimeStamp()
	serviceRegistry.Services.Range(func(k, service interface{}) bool {
		s := service.(Service)
		willBeEliminated := now-s.Timestamp > int64(serviceRegistry.HeartBeat)
		if willBeEliminated {
			serviceRegistry.Unregister(s.Name, s.Ip, s.Port)
		}
		return true
	})
}

func (serviceRegistry *ServiceRegistry) InitTick() {
	for range time.Tick(time.Duration(serviceRegistry.HeartBeat) * time.Second) {
		serviceRegistry.Cleanup()
	}
}
