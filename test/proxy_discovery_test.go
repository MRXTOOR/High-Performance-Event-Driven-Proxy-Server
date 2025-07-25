package test

import (
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/adminapi"
	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/config"
	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/core"
	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/discovery"
)

func TestDiscoveryDynamicBackends(t *testing.T) {
	backend1 := http.Server{Addr: ":18481", Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("d1"))
	})}
	backend2 := http.Server{Addr: ":18482", Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("d2"))
	})}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); backend1.ListenAndServe() }()
	go func() { defer wg.Done(); backend2.ListenAndServe() }()
	defer backend1.Close()
	defer backend2.Close()
	time.Sleep(200 * time.Millisecond)

	cfg := &config.Config{
		Listen:   config.ListenConfig{HTTP: ":18480"},
		Backends: []config.BackendConfig{},
	}
	proxyRef := adminapi.NewProxyRef(cfg)
	go proxyRef.Value.Load().(*core.ProxyServer).Start()
	time.Sleep(300 * time.Millisecond)

	// Сначала только backend1
	proxyRef.Value.Load().(*core.ProxyServer).SetBackends([]discovery.Backend{{Name: "d1", Address: "localhost", Port: 18481, Protocol: "http"}})
	time.Sleep(200 * time.Millisecond)
	resp, err := http.Get("http://localhost:18480/")
	if err != nil {
		t.Fatalf("Ошибка запроса: %v", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if string(body) != "d1" {
		t.Errorf("Ожидался ответ d1, а не %s", string(body))
	}

	// Теперь добавляем backend2
	proxyRef.Value.Load().(*core.ProxyServer).SetBackends([]discovery.Backend{{Name: "d1", Address: "localhost", Port: 18481, Protocol: "http"}, {Name: "d2", Address: "localhost", Port: 18482, Protocol: "http"}})
	foundD2 := false
	for poll := 0; poll < 10; poll++ {
		resp, err := http.Get("http://localhost:18480/")
		if err != nil {
			t.Fatalf("Ошибка запроса: %v", err)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if string(body) == "d2" {
			foundD2 = true
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	if !foundD2 {
		t.Errorf("Discovery не подхватил новый backend2")
	}

	// Теперь убираем backend1, оставляем только backend2
	proxyRef.Value.Load().(*core.ProxyServer).SetBackends([]discovery.Backend{{Name: "d2", Address: "localhost", Port: 18482, Protocol: "http"}})
	time.Sleep(200 * time.Millisecond)
	resp, err = http.Get("http://localhost:18480/")
	if err != nil {
		t.Fatalf("Ошибка запроса: %v", err)
	}
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if string(body) != "d2" {
		t.Errorf("Ожидался ответ d2, а не %s", string(body))
	}
}
