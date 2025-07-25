package discovery

import (
	"log"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

type Backend struct {
	Name     string
	Address  string
	Port     int
	Protocol string
}

type UpdateFunc func([]Backend)

func WatchConsulServices(consulAddr, serviceName, protocol string, update UpdateFunc) {
	cfg := consulapi.DefaultConfig()
	cfg.Address = consulAddr
	client, err := consulapi.NewClient(cfg)
	if err != nil {
		log.Printf("[DISCOVERY] Ошибка подключения к Consul: %v", err)
		return
	}
	var lastIndex uint64
	for {
		services, meta, err := client.Health().Service(serviceName, "", true, &consulapi.QueryOptions{
			WaitIndex: lastIndex,
			WaitTime:  30 * time.Second,
		})
		if err != nil {
			log.Printf("[DISCOVERY] Ошибка запроса к Consul: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}
		lastIndex = meta.LastIndex
		var backends []Backend
		for _, entry := range services {
			addr := entry.Service.Address
			if addr == "" {
				addr = entry.Node.Address
			}
			backends = append(backends, Backend{
				Name:     entry.Service.Service,
				Address:  addr,
				Port:     entry.Service.Port,
				Protocol: protocol,
			})
		}
		update(backends)
	}
}
