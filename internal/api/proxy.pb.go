
package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Empty struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Empty) Reset() {
	*x = Empty{}
	mi := &file_internal_api_proxy_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_proxy_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*Empty) Descriptor() ([]byte, []int) {
	return file_internal_api_proxy_proto_rawDescGZIP(), []int{0}
}

type ConfigResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ConfigYaml    string                 `protobuf:"bytes,1,opt,name=config_yaml,json=configYaml,proto3" json:"config_yaml,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ConfigResponse) Reset() {
	*x = ConfigResponse{}
	mi := &file_internal_api_proxy_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ConfigResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigResponse) ProtoMessage() {}

func (x *ConfigResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_proxy_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ConfigResponse) Descriptor() ([]byte, []int) {
	return file_internal_api_proxy_proto_rawDescGZIP(), []int{1}
}

func (x *ConfigResponse) GetConfigYaml() string {
	if x != nil {
		return x.ConfigYaml
	}
	return ""
}

type ReloadResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Error         string                 `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReloadResponse) Reset() {
	*x = ReloadResponse{}
	mi := &file_internal_api_proxy_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReloadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReloadResponse) ProtoMessage() {}

func (x *ReloadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_proxy_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ReloadResponse) Descriptor() ([]byte, []int) {
	return file_internal_api_proxy_proto_rawDescGZIP(), []int{2}
}

func (x *ReloadResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *ReloadResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type BackendsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Backends      []*Backend             `protobuf:"bytes,1,rep,name=backends,proto3" json:"backends,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BackendsResponse) Reset() {
	*x = BackendsResponse{}
	mi := &file_internal_api_proxy_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BackendsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BackendsResponse) ProtoMessage() {}

func (x *BackendsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_proxy_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*BackendsResponse) Descriptor() ([]byte, []int) {
	return file_internal_api_proxy_proto_rawDescGZIP(), []int{3}
}

func (x *BackendsResponse) GetBackends() []*Backend {
	if x != nil {
		return x.Backends
	}
	return nil
}

type Backend struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Address       string                 `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	Protocol      string                 `protobuf:"bytes,3,opt,name=protocol,proto3" json:"protocol,omitempty"`
	Alive         bool                   `protobuf:"varint,4,opt,name=alive,proto3" json:"alive,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Backend) Reset() {
	*x = Backend{}
	mi := &file_internal_api_proxy_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Backend) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Backend) ProtoMessage() {}

func (x *Backend) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_proxy_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*Backend) Descriptor() ([]byte, []int) {
	return file_internal_api_proxy_proto_rawDescGZIP(), []int{4}
}

func (x *Backend) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Backend) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *Backend) GetProtocol() string {
	if x != nil {
		return x.Protocol
	}
	return ""
}

func (x *Backend) GetAlive() bool {
	if x != nil {
		return x.Alive
	}
	return false
}

var File_internal_api_proxy_proto protoreflect.FileDescriptor

const file_internal_api_proxy_proto_rawDesc = "" +
	"\n" +
	"\x18internal/api/proxy.proto\x12\x05proxy\"\a\n" +
	"\x05Empty\"1\n" +
	"\x0eConfigResponse\x12\x1f\n" +
	"\vconfig_yaml\x18\x01 \x01(\tR\n" +
	"configYaml\"@\n" +
	"\x0eReloadResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\x12\x14\n" +
	"\x05error\x18\x02 \x01(\tR\x05error\">\n" +
	"\x10BackendsResponse\x12*\n" +
	"\bbackends\x18\x01 \x03(\v2\x0e.proxy.BackendR\bbackends\"i\n" +
	"\aBackend\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12\x18\n" +
	"\aaddress\x18\x02 \x01(\tR\aaddress\x12\x1a\n" +
	"\bprotocol\x18\x03 \x01(\tR\bprotocol\x12\x14\n" +
	"\x05alive\x18\x04 \x01(\bR\x05alive2\xa9\x01\n" +
	"\n" +
	"ProxyAdmin\x120\n" +
	"\tGetConfig\x12\f.proxy.Empty\x1a\x15.proxy.ConfigResponse\x123\n" +
	"\fReloadConfig\x12\f.proxy.Empty\x1a\x15.proxy.ReloadResponse\x124\n" +
	"\vGetBackends\x12\f.proxy.Empty\x1a\x17.proxy.BackendsResponseBPZNgithub.com/MRXTOOR/High-Performance-Event-Driven-Proxy-Server/internal/api;apib\x06proto3"

var (
	file_internal_api_proxy_proto_rawDescOnce sync.Once
	file_internal_api_proxy_proto_rawDescData []byte
)

func file_internal_api_proxy_proto_rawDescGZIP() []byte {
	file_internal_api_proxy_proto_rawDescOnce.Do(func() {
		file_internal_api_proxy_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_internal_api_proxy_proto_rawDesc), len(file_internal_api_proxy_proto_rawDesc)))
	})
	return file_internal_api_proxy_proto_rawDescData
}

var file_internal_api_proxy_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_internal_api_proxy_proto_goTypes = []any{
	(*Empty)(nil),            // 0: proxy.Empty
	(*ConfigResponse)(nil),   // 1: proxy.ConfigResponse
	(*ReloadResponse)(nil),   // 2: proxy.ReloadResponse
	(*BackendsResponse)(nil), // 3: proxy.BackendsResponse
	(*Backend)(nil),          // 4: proxy.Backend
}
var file_internal_api_proxy_proto_depIdxs = []int32{
	4, // 0: proxy.BackendsResponse.backends:type_name -> proxy.Backend
	0, // 1: proxy.ProxyAdmin.GetConfig:input_type -> proxy.Empty
	0, // 2: proxy.ProxyAdmin.ReloadConfig:input_type -> proxy.Empty
	0, // 3: proxy.ProxyAdmin.GetBackends:input_type -> proxy.Empty
	1, // 4: proxy.ProxyAdmin.GetConfig:output_type -> proxy.ConfigResponse
	2, // 5: proxy.ProxyAdmin.ReloadConfig:output_type -> proxy.ReloadResponse
	3, // 6: proxy.ProxyAdmin.GetBackends:output_type -> proxy.BackendsResponse
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_internal_api_proxy_proto_init() }
func file_internal_api_proxy_proto_init() {
	if File_internal_api_proxy_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_internal_api_proxy_proto_rawDesc), len(file_internal_api_proxy_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_api_proxy_proto_goTypes,
		DependencyIndexes: file_internal_api_proxy_proto_depIdxs,
		MessageInfos:      file_internal_api_proxy_proto_msgTypes,
	}.Build()
	File_internal_api_proxy_proto = out.File
	file_internal_api_proxy_proto_goTypes = nil
	file_internal_api_proxy_proto_depIdxs = nil
}
