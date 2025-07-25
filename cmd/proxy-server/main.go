package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/adminapi"
	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/config"
	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/core"
	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/discovery"
	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/metrics"
)

func main() {
	cfgPath := "configs/example.yaml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", err)
	}
	log.Println("Конфиг успешно загружен")

	var proxy atomic.Value
	proxy.Store(core.NewProxyServer(cfg))
	go func() {
		proxy.Load().(*core.ProxyServer).Start()
	}()

	go adminapi.RunGRPCServer(":9090", cfgPath, &proxy)

	if cfg.Discovery.Type == "consul" {
		go discovery.WatchConsulServices(
			cfg.Discovery.Address,
			cfg.Discovery.ServiceName,
			cfg.Discovery.Protocol,
			func(backends []discovery.Backend) {
				proxy.Load().(*core.ProxyServer).SetBackends(backends)
			},
		)
	}

	config.StartConfigWatcher(cfgPath, func(path string) bool {
		newCfg, err := config.LoadConfig(path)
		if err != nil {
			log.Printf("[HOTRELOAD] Ошибка загрузки конфига: %v", err)
			return false
		}
		log.Printf("[HOTRELOAD] Конфиг успешно перезагружен")
		proxy.Store(core.NewProxyServer(newCfg))
		go func() {
			proxy.Load().(*core.ProxyServer).Start()
		}()
		return true
	})

	if shutdown, err := metrics.InitTracing("proxy-server"); err != nil {
		log.Fatalf("Ошибка инициализации трассировки: %v", err)
	} else {
		defer shutdown(context.Background())
	}

	metrics.Init()
	go func() {
		log.Println("[METRICS] Экспорт /metrics на :9100")
		http.Handle("/metrics", metrics.Handler())
		http.ListenAndServe(":9100", nil)
	}()

	select {}
}
