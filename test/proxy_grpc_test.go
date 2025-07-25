package test

import (
	"context"
	"io/ioutil"
	"net"
	"os"
	"testing"
	"time"

	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/adminapi"
	pb "github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/api"
	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/config"
	"google.golang.org/grpc"
)

func TestGRPCAPI(t *testing.T) {
	cfgPath := "test_grpcapi.yaml"
	cfg := `listen:
  http: ":18580"
backends:
  - name: "b1"
    address: "localhost:18581"
    protocol: "http"
`
	if err := ioutil.WriteFile(cfgPath, []byte(cfg), 0644); err != nil {
		t.Fatalf("Ошибка записи конфига: %v", err)
	}
	defer os.Remove(cfgPath)

	proxyRef := newProxyRef(t, cfgPath)
	ln, err := net.Listen("tcp", ":19090")
	if err != nil {
		t.Fatalf("Ошибка listen: %v", err)
	}
	go adminapi.RunGRPCServerOnListener(ln, cfgPath, proxyRef)
	time.Sleep(300 * time.Millisecond)

	conn, err := grpc.Dial("localhost:19090", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
	if err != nil {
		t.Fatalf("Ошибка подключения к gRPC: %v", err)
	}
	defer conn.Close()
	client := pb.NewProxyAdminClient(conn)

	// GetConfig
	resp, err := client.GetConfig(context.Background(), &pb.Empty{})
	if err != nil || resp == nil || resp.ConfigYaml == "" {
		t.Errorf("GetConfig не вернул конфиг: %v", err)
	}

	// GetBackends
	bResp, err := client.GetBackends(context.Background(), &pb.Empty{})
	if err != nil || len(bResp.Backends) == 0 {
		t.Errorf("GetBackends не вернул backend'ы: %v", err)
	}

	// ReloadConfig (без изменений, просто проверяем что не падает)
	rResp, err := client.ReloadConfig(context.Background(), &pb.Empty{})
	if err != nil || !rResp.Success {
		t.Errorf("ReloadConfig не сработал: %v", err)
	}
}

func newProxyRef(t *testing.T, cfgPath string) *adminapi.ProxyRef {
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("Ошибка загрузки конфига: %v", err)
	}
	return adminapi.NewProxyRef(cfg)
}
