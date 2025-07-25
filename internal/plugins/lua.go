package plugins

import (
	"log"
	"net/http"

	lua "github.com/yuin/gopher-lua"
)

type LuaPlugin struct {
	ScriptPath string
}

func NewLuaPlugin(path string) *LuaPlugin {
	return &LuaPlugin{ScriptPath: path}
}

func (p *LuaPlugin) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		L := lua.NewState()
		defer L.Close()

		L.SetGlobal("method", lua.LString(r.Method))
		L.SetGlobal("url", lua.LString(r.URL.String()))
		L.SetGlobal("remote_addr", lua.LString(r.RemoteAddr))

		if err := L.DoFile(p.ScriptPath); err != nil {
			log.Printf("[LUA] Ошибка выполнения %s: %v", p.ScriptPath, err)
		}
		next.ServeHTTP(w, r)
	})
}
