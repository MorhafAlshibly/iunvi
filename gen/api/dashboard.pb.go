// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: api/dashboard.proto

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

type CreateDashboardRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ModelId       string                 `protobuf:"bytes,1,opt,name=modelId,proto3" json:"modelId,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Definition    string                 `protobuf:"bytes,3,opt,name=definition,proto3" json:"definition,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateDashboardRequest) Reset() {
	*x = CreateDashboardRequest{}
	mi := &file_api_dashboard_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateDashboardRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateDashboardRequest) ProtoMessage() {}

func (x *CreateDashboardRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_dashboard_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateDashboardRequest.ProtoReflect.Descriptor instead.
func (*CreateDashboardRequest) Descriptor() ([]byte, []int) {
	return file_api_dashboard_proto_rawDescGZIP(), []int{0}
}

func (x *CreateDashboardRequest) GetModelId() string {
	if x != nil {
		return x.ModelId
	}
	return ""
}

func (x *CreateDashboardRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateDashboardRequest) GetDefinition() string {
	if x != nil {
		return x.Definition
	}
	return ""
}

type CreateDashboardResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateDashboardResponse) Reset() {
	*x = CreateDashboardResponse{}
	mi := &file_api_dashboard_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateDashboardResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateDashboardResponse) ProtoMessage() {}

func (x *CreateDashboardResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_dashboard_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateDashboardResponse.ProtoReflect.Descriptor instead.
func (*CreateDashboardResponse) Descriptor() ([]byte, []int) {
	return file_api_dashboard_proto_rawDescGZIP(), []int{1}
}

func (x *CreateDashboardResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetDashboardsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	WorkspaceId   string                 `protobuf:"bytes,1,opt,name=workspaceId,proto3" json:"workspaceId,omitempty"`
	ModelId       *string                `protobuf:"bytes,2,opt,name=modelId,proto3,oneof" json:"modelId,omitempty"`
	ModelRunId    *string                `protobuf:"bytes,3,opt,name=modelRunId,proto3,oneof" json:"modelRunId,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetDashboardsRequest) Reset() {
	*x = GetDashboardsRequest{}
	mi := &file_api_dashboard_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetDashboardsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDashboardsRequest) ProtoMessage() {}

func (x *GetDashboardsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_dashboard_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDashboardsRequest.ProtoReflect.Descriptor instead.
func (*GetDashboardsRequest) Descriptor() ([]byte, []int) {
	return file_api_dashboard_proto_rawDescGZIP(), []int{2}
}

func (x *GetDashboardsRequest) GetWorkspaceId() string {
	if x != nil {
		return x.WorkspaceId
	}
	return ""
}

func (x *GetDashboardsRequest) GetModelId() string {
	if x != nil && x.ModelId != nil {
		return *x.ModelId
	}
	return ""
}

func (x *GetDashboardsRequest) GetModelRunId() string {
	if x != nil && x.ModelRunId != nil {
		return *x.ModelRunId
	}
	return ""
}

type GetDashboardsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Dashboards    []*Dashboard           `protobuf:"bytes,1,rep,name=dashboards,proto3" json:"dashboards,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetDashboardsResponse) Reset() {
	*x = GetDashboardsResponse{}
	mi := &file_api_dashboard_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetDashboardsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDashboardsResponse) ProtoMessage() {}

func (x *GetDashboardsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_dashboard_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDashboardsResponse.ProtoReflect.Descriptor instead.
func (*GetDashboardsResponse) Descriptor() ([]byte, []int) {
	return file_api_dashboard_proto_rawDescGZIP(), []int{3}
}

func (x *GetDashboardsResponse) GetDashboards() []*Dashboard {
	if x != nil {
		return x.Dashboards
	}
	return nil
}

type GetDashboardRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetDashboardRequest) Reset() {
	*x = GetDashboardRequest{}
	mi := &file_api_dashboard_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetDashboardRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDashboardRequest) ProtoMessage() {}

func (x *GetDashboardRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_dashboard_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDashboardRequest.ProtoReflect.Descriptor instead.
func (*GetDashboardRequest) Descriptor() ([]byte, []int) {
	return file_api_dashboard_proto_rawDescGZIP(), []int{4}
}

func (x *GetDashboardRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetDashboardResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Dashboard     *Dashboard             `protobuf:"bytes,1,opt,name=dashboard,proto3" json:"dashboard,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetDashboardResponse) Reset() {
	*x = GetDashboardResponse{}
	mi := &file_api_dashboard_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetDashboardResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDashboardResponse) ProtoMessage() {}

func (x *GetDashboardResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_dashboard_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDashboardResponse.ProtoReflect.Descriptor instead.
func (*GetDashboardResponse) Descriptor() ([]byte, []int) {
	return file_api_dashboard_proto_rawDescGZIP(), []int{5}
}

func (x *GetDashboardResponse) GetDashboard() *Dashboard {
	if x != nil {
		return x.Dashboard
	}
	return nil
}

type GetDashboardMarkdownResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Markdown      string                 `protobuf:"bytes,1,opt,name=markdown,proto3" json:"markdown,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetDashboardMarkdownResponse) Reset() {
	*x = GetDashboardMarkdownResponse{}
	mi := &file_api_dashboard_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetDashboardMarkdownResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDashboardMarkdownResponse) ProtoMessage() {}

func (x *GetDashboardMarkdownResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_dashboard_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDashboardMarkdownResponse.ProtoReflect.Descriptor instead.
func (*GetDashboardMarkdownResponse) Descriptor() ([]byte, []int) {
	return file_api_dashboard_proto_rawDescGZIP(), []int{6}
}

func (x *GetDashboardMarkdownResponse) GetMarkdown() string {
	if x != nil {
		return x.Markdown
	}
	return ""
}

type GetModelRunDashboardRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ModelRunId    string                 `protobuf:"bytes,1,opt,name=modelRunId,proto3" json:"modelRunId,omitempty"`
	DashboardId   string                 `protobuf:"bytes,2,opt,name=dashboardId,proto3" json:"dashboardId,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetModelRunDashboardRequest) Reset() {
	*x = GetModelRunDashboardRequest{}
	mi := &file_api_dashboard_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetModelRunDashboardRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetModelRunDashboardRequest) ProtoMessage() {}

func (x *GetModelRunDashboardRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_dashboard_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetModelRunDashboardRequest.ProtoReflect.Descriptor instead.
func (*GetModelRunDashboardRequest) Descriptor() ([]byte, []int) {
	return file_api_dashboard_proto_rawDescGZIP(), []int{7}
}

func (x *GetModelRunDashboardRequest) GetModelRunId() string {
	if x != nil {
		return x.ModelRunId
	}
	return ""
}

func (x *GetModelRunDashboardRequest) GetDashboardId() string {
	if x != nil {
		return x.DashboardId
	}
	return ""
}

type GetModelRunDashboardResponse struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	DashboardSasUrl string                 `protobuf:"bytes,1,opt,name=dashboardSasUrl,proto3" json:"dashboardSasUrl,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *GetModelRunDashboardResponse) Reset() {
	*x = GetModelRunDashboardResponse{}
	mi := &file_api_dashboard_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetModelRunDashboardResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetModelRunDashboardResponse) ProtoMessage() {}

func (x *GetModelRunDashboardResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_dashboard_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetModelRunDashboardResponse.ProtoReflect.Descriptor instead.
func (*GetModelRunDashboardResponse) Descriptor() ([]byte, []int) {
	return file_api_dashboard_proto_rawDescGZIP(), []int{8}
}

func (x *GetModelRunDashboardResponse) GetDashboardSasUrl() string {
	if x != nil {
		return x.DashboardSasUrl
	}
	return ""
}

type Dashboard struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ModelId       string                 `protobuf:"bytes,2,opt,name=modelId,proto3" json:"modelId,omitempty"`
	Name          string                 `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Dashboard) Reset() {
	*x = Dashboard{}
	mi := &file_api_dashboard_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Dashboard) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Dashboard) ProtoMessage() {}

func (x *Dashboard) ProtoReflect() protoreflect.Message {
	mi := &file_api_dashboard_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Dashboard.ProtoReflect.Descriptor instead.
func (*Dashboard) Descriptor() ([]byte, []int) {
	return file_api_dashboard_proto_rawDescGZIP(), []int{9}
}

func (x *Dashboard) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Dashboard) GetModelId() string {
	if x != nil {
		return x.ModelId
	}
	return ""
}

func (x *Dashboard) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_api_dashboard_proto protoreflect.FileDescriptor

const file_api_dashboard_proto_rawDesc = "" +
	"\n" +
	"\x13api/dashboard.proto\x12\tdashboard\"f\n" +
	"\x16CreateDashboardRequest\x12\x18\n" +
	"\amodelId\x18\x01 \x01(\tR\amodelId\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12\x1e\n" +
	"\n" +
	"definition\x18\x03 \x01(\tR\n" +
	"definition\")\n" +
	"\x17CreateDashboardResponse\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\"\x97\x01\n" +
	"\x14GetDashboardsRequest\x12 \n" +
	"\vworkspaceId\x18\x01 \x01(\tR\vworkspaceId\x12\x1d\n" +
	"\amodelId\x18\x02 \x01(\tH\x00R\amodelId\x88\x01\x01\x12#\n" +
	"\n" +
	"modelRunId\x18\x03 \x01(\tH\x01R\n" +
	"modelRunId\x88\x01\x01B\n" +
	"\n" +
	"\b_modelIdB\r\n" +
	"\v_modelRunId\"M\n" +
	"\x15GetDashboardsResponse\x124\n" +
	"\n" +
	"dashboards\x18\x01 \x03(\v2\x14.dashboard.DashboardR\n" +
	"dashboards\"%\n" +
	"\x13GetDashboardRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\"J\n" +
	"\x14GetDashboardResponse\x122\n" +
	"\tdashboard\x18\x01 \x01(\v2\x14.dashboard.DashboardR\tdashboard\":\n" +
	"\x1cGetDashboardMarkdownResponse\x12\x1a\n" +
	"\bmarkdown\x18\x01 \x01(\tR\bmarkdown\"_\n" +
	"\x1bGetModelRunDashboardRequest\x12\x1e\n" +
	"\n" +
	"modelRunId\x18\x01 \x01(\tR\n" +
	"modelRunId\x12 \n" +
	"\vdashboardId\x18\x02 \x01(\tR\vdashboardId\"H\n" +
	"\x1cGetModelRunDashboardResponse\x12(\n" +
	"\x0fdashboardSasUrl\x18\x01 \x01(\tR\x0fdashboardSasUrl\"I\n" +
	"\tDashboard\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x18\n" +
	"\amodelId\x18\x02 \x01(\tR\amodelId\x12\x12\n" +
	"\x04name\x18\x03 \x01(\tR\x04name2\xdb\x03\n" +
	"\x10DashboardService\x12X\n" +
	"\x0fCreateDashboard\x12!.dashboard.CreateDashboardRequest\x1a\".dashboard.CreateDashboardResponse\x12R\n" +
	"\rGetDashboards\x12\x1f.dashboard.GetDashboardsRequest\x1a .dashboard.GetDashboardsResponse\x12O\n" +
	"\fGetDashboard\x12\x1e.dashboard.GetDashboardRequest\x1a\x1f.dashboard.GetDashboardResponse\x12_\n" +
	"\x14GetDashboardMarkdown\x12\x1e.dashboard.GetDashboardRequest\x1a'.dashboard.GetDashboardMarkdownResponse\x12g\n" +
	"\x14GetModelRunDashboard\x12&.dashboard.GetModelRunDashboardRequest\x1a'.dashboard.GetModelRunDashboardResponseB)Z'github.com/MorhafAlshibly/iunvi/gen/apib\x06proto3"

var (
	file_api_dashboard_proto_rawDescOnce sync.Once
	file_api_dashboard_proto_rawDescData []byte
)

func file_api_dashboard_proto_rawDescGZIP() []byte {
	file_api_dashboard_proto_rawDescOnce.Do(func() {
		file_api_dashboard_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_dashboard_proto_rawDesc), len(file_api_dashboard_proto_rawDesc)))
	})
	return file_api_dashboard_proto_rawDescData
}

var file_api_dashboard_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_api_dashboard_proto_goTypes = []any{
	(*CreateDashboardRequest)(nil),       // 0: dashboard.CreateDashboardRequest
	(*CreateDashboardResponse)(nil),      // 1: dashboard.CreateDashboardResponse
	(*GetDashboardsRequest)(nil),         // 2: dashboard.GetDashboardsRequest
	(*GetDashboardsResponse)(nil),        // 3: dashboard.GetDashboardsResponse
	(*GetDashboardRequest)(nil),          // 4: dashboard.GetDashboardRequest
	(*GetDashboardResponse)(nil),         // 5: dashboard.GetDashboardResponse
	(*GetDashboardMarkdownResponse)(nil), // 6: dashboard.GetDashboardMarkdownResponse
	(*GetModelRunDashboardRequest)(nil),  // 7: dashboard.GetModelRunDashboardRequest
	(*GetModelRunDashboardResponse)(nil), // 8: dashboard.GetModelRunDashboardResponse
	(*Dashboard)(nil),                    // 9: dashboard.Dashboard
}
var file_api_dashboard_proto_depIdxs = []int32{
	9, // 0: dashboard.GetDashboardsResponse.dashboards:type_name -> dashboard.Dashboard
	9, // 1: dashboard.GetDashboardResponse.dashboard:type_name -> dashboard.Dashboard
	0, // 2: dashboard.DashboardService.CreateDashboard:input_type -> dashboard.CreateDashboardRequest
	2, // 3: dashboard.DashboardService.GetDashboards:input_type -> dashboard.GetDashboardsRequest
	4, // 4: dashboard.DashboardService.GetDashboard:input_type -> dashboard.GetDashboardRequest
	4, // 5: dashboard.DashboardService.GetDashboardMarkdown:input_type -> dashboard.GetDashboardRequest
	7, // 6: dashboard.DashboardService.GetModelRunDashboard:input_type -> dashboard.GetModelRunDashboardRequest
	1, // 7: dashboard.DashboardService.CreateDashboard:output_type -> dashboard.CreateDashboardResponse
	3, // 8: dashboard.DashboardService.GetDashboards:output_type -> dashboard.GetDashboardsResponse
	5, // 9: dashboard.DashboardService.GetDashboard:output_type -> dashboard.GetDashboardResponse
	6, // 10: dashboard.DashboardService.GetDashboardMarkdown:output_type -> dashboard.GetDashboardMarkdownResponse
	8, // 11: dashboard.DashboardService.GetModelRunDashboard:output_type -> dashboard.GetModelRunDashboardResponse
	7, // [7:12] is the sub-list for method output_type
	2, // [2:7] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_dashboard_proto_init() }
func file_api_dashboard_proto_init() {
	if File_api_dashboard_proto != nil {
		return
	}
	file_api_dashboard_proto_msgTypes[2].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_dashboard_proto_rawDesc), len(file_api_dashboard_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_dashboard_proto_goTypes,
		DependencyIndexes: file_api_dashboard_proto_depIdxs,
		MessageInfos:      file_api_dashboard_proto_msgTypes,
	}.Build()
	File_api_dashboard_proto = out.File
	file_api_dashboard_proto_goTypes = nil
	file_api_dashboard_proto_depIdxs = nil
}
