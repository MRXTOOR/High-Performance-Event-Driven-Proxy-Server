package core

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/config"
	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/discovery"
	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/metrics"
	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/plugins"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ProxyServer struct {
	Config  *config.Config
	mu      sync.RWMutex
	httpSrv *http.Server
	handler atomic.Value
}

func NewProxyServer(cfg *config.Config) *ProxyServer {
	p := &ProxyServer{Config: cfg}
	p.UpdateHandler(p.buildHandler())
	return p
}

func (p *ProxyServer) SetBackends(backends []discovery.Backend) {
	p.mu.Lock()
	defer p.mu.Unlock()
	var newBackends []config.BackendConfig
	for _, b := range backends {
		newBackends = append(newBackends, config.BackendConfig{
			Name:     b.Name,
			Address:  b.Address + ":" + itoa(b.Port),
			Protocol: b.Protocol,
		})
	}
	p.Config.Backends = newBackends
	p.UpdateHandler(p.buildHandler())
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}

func (p *ProxyServer) Start() {
	if p.Config.Listen.HTTP != "" {
		mux := http.NewServeMux()
		mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p.getHandler().ServeHTTP(w, r)
		}))
		p.httpSrv = &http.Server{Addr: p.Config.Listen.HTTP, Handler: mux}
		log.Printf("[HTTP] Сервер слушает на %s", p.Config.Listen.HTTP)
		if err := p.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[HTTP] Ошибка: %v", err)
		}
	}
	if p.Config.Listen.TCP != "" {
		go func() {
			ln, err := net.Listen("tcp", p.Config.Listen.TCP)
			if err != nil {
				log.Printf("[TCP] Ошибка: %v", err)
				return
			}
			log.Printf("[TCP] Сервер слушает на %s", p.Config.Listen.TCP)
			for {
				conn, err := ln.Accept()
				if err != nil {
					log.Printf("[TCP] Accept error: %v", err)
					continue
				}
				go func(c net.Conn) {
					defer c.Close()
					io.Copy(c, c)
				}(conn)
			}
		}()
	}
	if p.Config.Listen.UDP != "" {
		go func() {
			addr, err := net.ResolveUDPAddr("udp", p.Config.Listen.UDP)
			if err != nil {
				log.Printf("[UDP] Ошибка: %v", err)
				return
			}
			conn, err := net.ListenUDP("udp", addr)
			if err != nil {
				log.Printf("[UDP] Ошибка: %v", err)
				return
			}
			log.Printf("[UDP] Сервер слушает на %s", p.Config.Listen.UDP)
			buf := make([]byte, 2048)
			for {
				n, remote, err := conn.ReadFromUDP(buf)
				if err != nil {
					log.Printf("[UDP] Read error: %v", err)
					continue
				}
				_, err = conn.WriteToUDP(buf[:n], remote)
				if err != nil {
					log.Printf("[UDP] Write error: %v", err)
				}
			}
		}()
	}
	fmt.Println("Proxy server завершил работу.")
}

func (p *ProxyServer) Shutdown(ctx context.Context) error {
	if p.httpSrv != nil {
		return p.httpSrv.Shutdown(ctx)
	}
	return nil
}

func (p *ProxyServer) getHandler() http.Handler {
	h, _ := p.handler.Load().(http.Handler)
	return h
}

func (p *ProxyServer) UpdateHandler(newHandler http.Handler) {
	p.handler.Store(newHandler)
}

func (p *ProxyServer) httpBackendURLs() []*url.URL {
	var urls []*url.URL
	for _, b := range p.Config.Backends {
		if b.Protocol == "http" {
			u, err := url.Parse("http://" + b.Address)
			if err == nil {
				urls = append(urls, u)
			}
		}
	}
	return urls
}

func checkBackendAlive(u *url.URL) bool {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(u.String())
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode < 500
}

func (p *ProxyServer) BackendStatus() []struct {
	Name     string
	Address  string
	Protocol string
	Alive    bool
} {
	var result []struct {
		Name     string
		Address  string
		Protocol string
		Alive    bool
	}
	for _, b := range p.Config.Backends {
		alive := true
		if b.Protocol == "http" {
			alive = false
			for _, u := range p.httpBackendURLs() {
				if u.Host == b.Address {
					alive = true
				}
			}
		}
		result = append(result, struct {
			Name     string
			Address  string
			Protocol string
			Alive    bool
		}{
			Name:     b.Name,
			Address:  b.Address,
			Protocol: b.Protocol,
			Alive:    alive,
		})
	}
	return result
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (p *ProxyServer) buildHandler() http.Handler {
	tracer := otel.Tracer("proxy-server")
	backends := p.httpBackendURLs()
	if len(backends) == 0 {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, span := tracer.Start(r.Context(), "http_request",
				trace.WithAttributes(
					attribute.String("http.method", r.Method),
					attribute.String("http.target", r.URL.Path),
				),
			)
			defer span.End()
			start := time.Now()
			msg := fmt.Sprintf("Echo HTTP: %s %s\n", r.Method, r.URL.Path)
			w.Write([]byte(msg))
			metrics.ObserveRequest(r.Method, r.URL.Path, 200, time.Since(start))
		})
	}

	var rrIndex uint64
	proxies := make([]*httputil.ReverseProxy, len(backends))
	alive := make([]bool, len(backends))
	for i, u := range backends {
		proxies[i] = httputil.NewSingleHostReverseProxy(u)
		alive[i] = true
	}
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			for i, u := range backends {
				alive[i] = checkBackendAlive(u)
			}
		}
	}()

	var handler http.Handler
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "http_request",
			trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.target", r.URL.Path),
			),
		)
		defer span.End()
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: 200}
		var liveProxies []*httputil.ReverseProxy
		for i, p := range proxies {
			if alive[i] {
				liveProxies = append(liveProxies, p)
			}
		}
		if len(liveProxies) == 0 {
			rw.WriteHeader(502)
			rw.Write([]byte("502 Bad Gateway: нет доступных backend'ов\n"))
			metrics.ObserveRequest(r.Method, r.URL.Path, 502, time.Since(start))
			return
		}
		i := atomic.AddUint64(&rrIndex, 1)
		proxy := liveProxies[int(i)%len(liveProxies)]
		proxy.ServeHTTP(rw, r.WithContext(ctx))
		metrics.ObserveRequest(r.Method, r.URL.Path, rw.status, time.Since(start))
	})

	for _, pl := range p.Config.Plugins {
		if pl.Type == "lua" {
			handler = plugins.NewLuaPlugin(pl.Path).Middleware(handler)
		}
		if pl.Type == "wasm" {
			handler = plugins.NewWasmPlugin(pl.Path).Middleware(handler)
		}
	}

	return handler
}

func (p *ProxyServer) BuildHandler() http.Handler {
	return p.buildHandler()
}
