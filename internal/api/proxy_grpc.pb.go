
package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const _ = grpc.SupportPackageIsVersion9

const (
	ProxyAdmin_GetConfig_FullMethodName    = "/proxy.ProxyAdmin/GetConfig"
	ProxyAdmin_ReloadConfig_FullMethodName = "/proxy.ProxyAdmin/ReloadConfig"
	ProxyAdmin_GetBackends_FullMethodName  = "/proxy.ProxyAdmin/GetBackends"
)

type ProxyAdminClient interface {
	GetConfig(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ConfigResponse, error)
	ReloadConfig(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ReloadResponse, error)
	GetBackends(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*BackendsResponse, error)
}

type proxyAdminClient struct {
	cc grpc.ClientConnInterface
}

func NewProxyAdminClient(cc grpc.ClientConnInterface) ProxyAdminClient {
	return &proxyAdminClient{cc}
}

func (c *proxyAdminClient) GetConfig(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ConfigResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ConfigResponse)
	err := c.cc.Invoke(ctx, ProxyAdmin_GetConfig_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyAdminClient) ReloadConfig(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ReloadResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReloadResponse)
	err := c.cc.Invoke(ctx, ProxyAdmin_ReloadConfig_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyAdminClient) GetBackends(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*BackendsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BackendsResponse)
	err := c.cc.Invoke(ctx, ProxyAdmin_GetBackends_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type ProxyAdminServer interface {
	GetConfig(context.Context, *Empty) (*ConfigResponse, error)
	ReloadConfig(context.Context, *Empty) (*ReloadResponse, error)
	GetBackends(context.Context, *Empty) (*BackendsResponse, error)
	mustEmbedUnimplementedProxyAdminServer()
}

type UnimplementedProxyAdminServer struct{}

func (UnimplementedProxyAdminServer) GetConfig(context.Context, *Empty) (*ConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConfig not implemented")
}
func (UnimplementedProxyAdminServer) ReloadConfig(context.Context, *Empty) (*ReloadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReloadConfig not implemented")
}
func (UnimplementedProxyAdminServer) GetBackends(context.Context, *Empty) (*BackendsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBackends not implemented")
}
func (UnimplementedProxyAdminServer) mustEmbedUnimplementedProxyAdminServer() {}
func (UnimplementedProxyAdminServer) testEmbeddedByValue()                    {}

type UnsafeProxyAdminServer interface {
	mustEmbedUnimplementedProxyAdminServer()
}

func RegisterProxyAdminServer(s grpc.ServiceRegistrar, srv ProxyAdminServer) {
	// If the following call pancis, it indicates UnimplementedProxyAdminServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ProxyAdmin_ServiceDesc, srv)
}

func _ProxyAdmin_GetConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyAdminServer).GetConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProxyAdmin_GetConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyAdminServer).GetConfig(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProxyAdmin_ReloadConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyAdminServer).ReloadConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProxyAdmin_ReloadConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyAdminServer).ReloadConfig(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProxyAdmin_GetBackends_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyAdminServer).GetBackends(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProxyAdmin_GetBackends_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyAdminServer).GetBackends(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var ProxyAdmin_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proxy.ProxyAdmin",
	HandlerType: (*ProxyAdminServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetConfig",
			Handler:    _ProxyAdmin_GetConfig_Handler,
		},
		{
			MethodName: "ReloadConfig",
			Handler:    _ProxyAdmin_ReloadConfig_Handler,
		},
		{
			MethodName: "GetBackends",
			Handler:    _ProxyAdmin_GetBackends_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/api/proxy.proto",
}
