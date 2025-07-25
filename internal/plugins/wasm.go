package plugins

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/tetratelabs/wazero"
)

type WasmPlugin struct {
	ModulePath string
}

func NewWasmPlugin(path string) *WasmPlugin {
	return &WasmPlugin{ModulePath: path}
}

func (p *WasmPlugin) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		runtime := wazero.NewRuntime(context.Background())
		defer runtime.Close(context.Background())

		wasmBytes, err := os.ReadFile(p.ModulePath)
		if err != nil {
			log.Printf("[WASM] Ошибка чтения %s: %v", p.ModulePath, err)
			next.ServeHTTP(w, r)
			return
		}

		mod, err := runtime.Instantiate(context.Background(), wasmBytes)
		if err != nil {
			log.Printf("[WASM] Ошибка запуска %s: %v", p.ModulePath, err)
			next.ServeHTTP(w, r)
			return
		}
		defer mod.Close(context.Background())

		// Пример: вызываем функцию "on_request"
		onReq := mod.ExportedFunction("on_request")
		if onReq != nil {
			_, err := onReq.Call(context.Background())
			if err != nil {
				log.Printf("[WASM] Ошибка вызова on_request: %v", err)
			}
		}

		next.ServeHTTP(w, r)
	})
}
