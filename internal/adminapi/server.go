package adminapi

import (
	context "context"
	"io/ioutil"
	"log"
	"net"
	"sync/atomic"

	pb "github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/api"

	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/config"
	"github.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)


type ProxyAdminService struct {
	pb.UnimplementedProxyAdminServer
	ConfigPath string
	ProxyRef   *atomic.Value // *core.ProxyServer
}

func NewProxyAdminService(configPath string, proxyRef *atomic.Value) *ProxyAdminService {
	return &ProxyAdminService{
		ConfigPath: configPath,
		ProxyRef:   proxyRef,
	}
}

func (s *ProxyAdminService) GetConfig(ctx context.Context, req *pb.Empty) (*pb.ConfigResponse, error) {
	data, err := ioutil.ReadFile(s.ConfigPath)
	if err != nil {
		return nil, err
	}
	return &pb.ConfigResponse{ConfigYaml: string(data)}, nil
}

func (s *ProxyAdminService) ReloadConfig(ctx context.Context, req *pb.Empty) (*pb.ReloadResponse, error) {
	cfg, err := config.LoadConfig(s.ConfigPath)
	if err != nil {
		return &pb.ReloadResponse{Success: false, Error: err.Error()}, nil
	}
	s.ProxyRef.Store(core.NewProxyServer(cfg))
	go func() { s.ProxyRef.Load().(*core.ProxyServer).Start() }()
	return &pb.ReloadResponse{Success: true}, nil
}

func (s *ProxyAdminService) GetBackends(ctx context.Context, req *pb.Empty) (*pb.BackendsResponse, error) {
	proxy := s.ProxyRef.Load().(*core.ProxyServer)
	backends := proxy.BackendStatus()
	resp := &pb.BackendsResponse{}
	for _, b := range backends {
		resp.Backends = append(resp.Backends, &pb.Backend{
			Name:     b.Name,
			Address:  b.Address,
			Protocol: b.Protocol,
			Alive:    b.Alive,
		})
	}
	return resp, nil
}

type ProxyRef struct {
	Value atomic.Value // *core.ProxyServer
}

func NewProxyRef(cfg *config.Config) *ProxyRef {
	ref := &ProxyRef{}
	ref.Value.Store(core.NewProxyServer(cfg))
	return ref
}

func RunGRPCServer(addr, configPath string, proxyRef *atomic.Value) {
	grpcServer := grpc.NewServer()
	pb.RegisterProxyAdminServer(grpcServer, NewProxyAdminService(configPath, proxyRef))
	reflection.Register(grpcServer)
	l, err := core.Listen(addr)
	if err != nil {
		log.Fatalf("[gRPC] Ошибка listen: %v", err)
	}
	log.Printf("[gRPC] Сервер слушает на %s", addr)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("[gRPC] Ошибка serve: %v", err)
	}
}

func RunGRPCServerOnListener(ln net.Listener, cfgPath string, proxyRef *ProxyRef) {
	grpcServer := grpc.NewServer()
	pb.RegisterProxyAdminServer(grpcServer, NewProxyAdminService(cfgPath, &proxyRef.Value))
	reflection.Register(grpcServer)
	grpcServer.Serve(ln)
}
