package test

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/config"
	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/core"
)

func TestHTTPWithWasmPlugin(t *testing.T) {
	wasmPath := "test_wasm_plugin.wasm"
	logFile := "wasm_plugin.log"
	_ = os.Remove(logFile)
	if _, err := os.Stat(wasmPath); err != nil {
		t.Skip("WASM-модуль не найден, пропускаем тест")
	}

	backend := http.Server{Addr: ":18281", Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})}
	go backend.ListenAndServe()
	defer backend.Close()
	time.Sleep(200 * time.Millisecond)

	cfg := &config.Config{
		Listen: config.ListenConfig{HTTP: ":18280"},
		Backends: []config.BackendConfig{
			{Name: "b1", Address: "localhost:18281", Protocol: "http"},
		},
		Plugins: []config.PluginConfig{
			{Type: "wasm", Path: wasmPath},
		},
	}
	proxy := core.NewProxyServer(cfg)
	go proxy.Start()
	time.Sleep(300 * time.Millisecond)

	resp, err := http.Get("http://localhost:18280/test-wasm")
	if err != nil {
		t.Fatalf("Ошибка запроса к прокси: %v", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if string(body) != "ok" {
		t.Errorf("Неверный ответ от backend: %s", string(body))
	}

	// Проверяем, что wasm-плагин сработал (файл создан и содержит строку)
	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Errorf("WASM-плагин не сработал: %v", err)
	}
	if string(data) == "" || string(data)[:4] != "wasm" {
		t.Errorf("WASM-плагин не записал ожидаемую строку: %s", string(data))
	}
}
