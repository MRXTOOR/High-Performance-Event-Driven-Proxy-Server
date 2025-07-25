package test

import (
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/adminapi"
	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/config"
	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/core"
)

func TestHotReloadConfig(t *testing.T) {
	cfgPath := "test_hotreload.yaml"
	backend1 := http.Server{Addr: ":18381", Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("b1"))
	})}
	backend2 := http.Server{Addr: ":18382", Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("b2"))
	})}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); backend1.ListenAndServe() }()
	go func() { defer wg.Done(); backend2.ListenAndServe() }()
	defer backend1.Close()
	defer backend2.Close()
	time.Sleep(200 * time.Millisecond)

	// Сохраняем начальный конфиг (только backend1)
	cfg := `listen:
  http: ":18380"
backends:
  - name: "b1"
    address: "localhost:18381"
    protocol: "http"
`
	if err := ioutil.WriteFile(cfgPath, []byte(cfg), 0644); err != nil {
		t.Fatalf("Ошибка записи конфига: %v", err)
	}
	defer os.Remove(cfgPath)

	proxyRef := adminapi.NewProxyRef(mustLoadConfig(t, cfgPath))
	go proxyRef.Value.Load().(*core.ProxyServer).Start()
	time.Sleep(300 * time.Millisecond)

	// Проверяем, что запросы идут только на backend1
	resp, err := http.Get("http://localhost:18380/")
	if err != nil {
		t.Fatalf("Ошибка запроса: %v", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if string(body) != "b1" {
		t.Errorf("Ожидался ответ b1, а не %s", string(body))
	}

	// Меняем конфиг: добавляем backend2
	cfg2 := `listen:
  http: ":18380"
backends:
  - name: "b1"
    address: "localhost:18381"
    protocol: "http"
  - name: "b2"
    address: "localhost:18382"
    protocol: "http"
`
	if err := ioutil.WriteFile(cfgPath, []byte(cfg2), 0644); err != nil {
		t.Fatalf("Ошибка обновления конфига: %v", err)
	}
	// Reload через proxyRef
	newCfg := mustLoadConfig(t, cfgPath)
	proxy := proxyRef.Value.Load().(*core.ProxyServer)
	proxy.Config = newCfg
	proxy.UpdateHandler(proxy.BuildHandler())

	// Polling: ждем, пока backend2 появится в ответах
	foundB2 := false
	for poll := 0; poll < 10; poll++ {
		resp, err := http.Get("http://localhost:18380/")
		if err != nil {
			t.Fatalf("Ошибка запроса: %v", err)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if string(body) == "b2" {
			foundB2 = true
			break
		}
		time.Sleep(300 * time.Millisecond)
	}
	if !foundB2 {
		t.Errorf("Hot reload не подхватил новый backend2")
	}
}

func mustLoadConfig(t *testing.T, path string) *config.Config {
	cfg, err := config.LoadConfig(path)
	if err != nil {
		t.Fatalf("Ошибка загрузки конфига: %v", err)
	}
	return cfg
}
