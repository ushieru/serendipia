package serviceregistry

import (
	"errors"
	"github.com/ushieru/serendipia/src/utils"
	"log"
	"sync"
	"time"
)

type ServiceRegistry struct {
	Services  sync.Map
	HeartBeat int64
	serviceRR map[string]int
}

func NewServiceRegistry(heartBeat int64) *ServiceRegistry {
	return &ServiceRegistry{
		HeartBeat: heartBeat,
		serviceRR: make(map[string]int),
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
		log.Println("[ServiceRegistry] Updated service:", name, "at", ip+":"+port)
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
	log.Println("[ServiceRegistry] Registered service:", name, "at", ip+":"+port)
	return key
}

func (serviceRegistry *ServiceRegistry) Get(name string) (*Service, error) {
	services := make([]Service, 0)
	serviceRegistry.Services.Range(func(k, value any) bool {
		service := value.(Service)
		if service.Name == name {
			services = append(services, service)
		}
		return true
	})
	if len(services) == 0 {
		return nil, errors.New("service not found")
	}
	if _, ok := serviceRegistry.serviceRR[name]; !ok {
		serviceRegistry.serviceRR[name] = 0
	}
	val := serviceRegistry.serviceRR[name]
	if val > len(services)-1 {
		serviceRegistry.serviceRR[name] = 0
		return &services[0], nil
	}
	serviceRegistry.serviceRR[name] = val + 1
	return &services[val], nil
}

func (serviceRegistry *ServiceRegistry) Unregister(name string, ip string, port string) {
	key := name + ip + port
	serviceRegistry.Services.Delete(key)
	log.Println("[ServiceRegistry] Unregister service:", name, "at", ip+":"+port)
}

func (serviceRegistry *ServiceRegistry) Cleanup() {
	now := utils.GetTimeStamp()
	serviceRegistry.Services.Range(func(k, service any) bool {
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
