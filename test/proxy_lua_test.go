package test

import (
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/config"
	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/core"
)

func TestHTTPWithLuaPlugin(t *testing.T) {
	// Создаем временный Lua-скрипт, который пишет в файл
	logFile := "test_lua_plugin.log"
	defer os.Remove(logFile)
	luaScript := `local f = io.open("` + logFile + `", "a"); f:write(method .. "," .. url .. "\n"); f:close()`
	pluginPath := "test_lua_plugin.lua"
	if err := ioutil.WriteFile(pluginPath, []byte(luaScript), 0644); err != nil {
		t.Fatalf("Ошибка создания lua-скрипта: %v", err)
	}
	defer os.Remove(pluginPath)

	// Запускаем тестовый backend
	backend := http.Server{Addr: ":18081", Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done(); backend.ListenAndServe() }()
	defer backend.Close()
	// Ждем старта
	time.Sleep(200 * time.Millisecond)

	cfg := &config.Config{
		Listen: config.ListenConfig{HTTP: ":18080"},
		Backends: []config.BackendConfig{
			{Name: "b1", Address: "localhost:18081", Protocol: "http"},
		},
		Plugins: []config.PluginConfig{
			{Type: "lua", Path: pluginPath},
		},
	}
	proxy := core.NewProxyServer(cfg)
	go proxy.Start()
	// Ждем старта прокси
	time.Sleep(300 * time.Millisecond)

	resp, err := http.Get("http://localhost:18080/test-lua")
	if err != nil {
		t.Fatalf("Ошибка запроса к прокси: %v", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if string(body) != "ok" {
		t.Errorf("Неверный ответ от backend: %s", string(body))
	}

	// Проверяем, что Lua-плагин сработал
	data, err := ioutil.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Lua-плагин не записал лог: %v", err)
	}
	if string(data) == "" || string(data)[:3] != "GET" {
		t.Errorf("Lua-плагин не сработал или неверный лог: %s", string(data))
	}
}
