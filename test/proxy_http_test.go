package test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/config"
	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/core"
)

func startTestBackend(t *testing.T, port int, id string, wg *sync.WaitGroup) *http.Server {
	h := http.NewServeMux()
	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("backend-" + id))
	})
	srv := &http.Server{Addr: ":" + itoa(port), Handler: h}
	wg.Add(1)
	go func() {
		defer wg.Done()
		srv.ListenAndServe()
	}()
	// Ждем старта
	time.Sleep(200 * time.Millisecond)
	return srv
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}

func TestHTTPProxy_BalancingAndHealth(t *testing.T) {
	var wg sync.WaitGroup
	backend1 := startTestBackend(t, 18181, "A", &wg)
	backend2 := startTestBackend(t, 18182, "B", &wg)

	cfg := &config.Config{
		Listen: config.ListenConfig{HTTP: ":18180"},
		Backends: []config.BackendConfig{
			{Name: "b1", Address: "localhost:18181", Protocol: "http"},
			{Name: "b2", Address: "localhost:18182", Protocol: "http"},
		},
	}
	proxy := core.NewProxyServer(cfg)
	go proxy.Start()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		backend1.Shutdown(ctx)
		backend2.Shutdown(ctx)
	}()
	// Ждем старта прокси
	time.Sleep(500 * time.Millisecond)

	// Проверяем балансировку (должны получить оба backend'а)
	results := make(map[string]int)
	for i := 0; i < 10; i++ {
		resp, err := http.Get("http://localhost:18180/")
		if err != nil {
			t.Fatalf("Ошибка запроса к прокси: %v", err)
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		results[string(body)]++
	}
	if len(results) < 2 {
		t.Errorf("Балансировка не работает, ответы только от одного backend: %v", results)
	}

	// Останавливаем один backend
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	backend2.Shutdown(ctx)

	// Polling для health-check
	var onlyA bool
	for poll := 0; poll < 10; poll++ {
		results = make(map[string]int)
		for i := 0; i < 5; i++ {
			resp, err := http.Get("http://localhost:18180/")
			if err != nil {
				t.Fatalf("Ошибка запроса к прокси: %v", err)
			}
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			results[string(body)]++
		}
		if len(results) == 1 && results["backend-A"] == 5 {
			onlyA = true
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if !onlyA {
		t.Errorf("Health-check не исключил недоступный backend: %v", results)
	}
}
