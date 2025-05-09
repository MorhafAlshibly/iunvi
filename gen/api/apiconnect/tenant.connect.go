// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: api/tenant.proto

package apiconnect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	api "github.com/MorhafAlshibly/iunvi/gen/api"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// TenantServiceName is the fully-qualified name of the TenantService service.
	TenantServiceName = "tenant.TenantService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// TenantServiceCreateWorkspaceProcedure is the fully-qualified name of the TenantService's
	// CreateWorkspace RPC.
	TenantServiceCreateWorkspaceProcedure = "/tenant.TenantService/CreateWorkspace"
	// TenantServiceGetWorkspacesProcedure is the fully-qualified name of the TenantService's
	// GetWorkspaces RPC.
	TenantServiceGetWorkspacesProcedure = "/tenant.TenantService/GetWorkspaces"
	// TenantServiceEditWorkspaceProcedure is the fully-qualified name of the TenantService's
	// EditWorkspace RPC.
	TenantServiceEditWorkspaceProcedure = "/tenant.TenantService/EditWorkspace"
	// TenantServiceGetUsersProcedure is the fully-qualified name of the TenantService's GetUsers RPC.
	TenantServiceGetUsersProcedure = "/tenant.TenantService/GetUsers"
	// TenantServiceGetUserWorkspaceAssignmentProcedure is the fully-qualified name of the
	// TenantService's GetUserWorkspaceAssignment RPC.
	TenantServiceGetUserWorkspaceAssignmentProcedure = "/tenant.TenantService/GetUserWorkspaceAssignment"
	// TenantServiceAssignUserToWorkspaceProcedure is the fully-qualified name of the TenantService's
	// AssignUserToWorkspace RPC.
	TenantServiceAssignUserToWorkspaceProcedure = "/tenant.TenantService/AssignUserToWorkspace"
)

// TenantServiceClient is a client for the tenant.TenantService service.
type TenantServiceClient interface {
	CreateWorkspace(context.Context, *connect.Request[api.CreateWorkspaceRequest]) (*connect.Response[api.CreateWorkspaceResponse], error)
	GetWorkspaces(context.Context, *connect.Request[api.GetWorkspacesRequest]) (*connect.Response[api.GetWorkspacesResponse], error)
	EditWorkspace(context.Context, *connect.Request[api.EditWorkspaceRequest]) (*connect.Response[api.EditWorkspaceResponse], error)
	GetUsers(context.Context, *connect.Request[api.GetUsersRequest]) (*connect.Response[api.GetUsersResponse], error)
	GetUserWorkspaceAssignment(context.Context, *connect.Request[api.GetUserWorkspaceAssignmentRequest]) (*connect.Response[api.GetUserWorkspaceAssignmentResponse], error)
	AssignUserToWorkspace(context.Context, *connect.Request[api.AssignUserToWorkspaceRequest]) (*connect.Response[api.AssignUserToWorkspaceResponse], error)
}

// NewTenantServiceClient constructs a client for the tenant.TenantService service. By default, it
// uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewTenantServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) TenantServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	tenantServiceMethods := api.File_api_tenant_proto.Services().ByName("TenantService").Methods()
	return &tenantServiceClient{
		createWorkspace: connect.NewClient[api.CreateWorkspaceRequest, api.CreateWorkspaceResponse](
			httpClient,
			baseURL+TenantServiceCreateWorkspaceProcedure,
			connect.WithSchema(tenantServiceMethods.ByName("CreateWorkspace")),
			connect.WithClientOptions(opts...),
		),
		getWorkspaces: connect.NewClient[api.GetWorkspacesRequest, api.GetWorkspacesResponse](
			httpClient,
			baseURL+TenantServiceGetWorkspacesProcedure,
			connect.WithSchema(tenantServiceMethods.ByName("GetWorkspaces")),
			connect.WithClientOptions(opts...),
		),
		editWorkspace: connect.NewClient[api.EditWorkspaceRequest, api.EditWorkspaceResponse](
			httpClient,
			baseURL+TenantServiceEditWorkspaceProcedure,
			connect.WithSchema(tenantServiceMethods.ByName("EditWorkspace")),
			connect.WithClientOptions(opts...),
		),
		getUsers: connect.NewClient[api.GetUsersRequest, api.GetUsersResponse](
			httpClient,
			baseURL+TenantServiceGetUsersProcedure,
			connect.WithSchema(tenantServiceMethods.ByName("GetUsers")),
			connect.WithClientOptions(opts...),
		),
		getUserWorkspaceAssignment: connect.NewClient[api.GetUserWorkspaceAssignmentRequest, api.GetUserWorkspaceAssignmentResponse](
			httpClient,
			baseURL+TenantServiceGetUserWorkspaceAssignmentProcedure,
			connect.WithSchema(tenantServiceMethods.ByName("GetUserWorkspaceAssignment")),
			connect.WithClientOptions(opts...),
		),
		assignUserToWorkspace: connect.NewClient[api.AssignUserToWorkspaceRequest, api.AssignUserToWorkspaceResponse](
			httpClient,
			baseURL+TenantServiceAssignUserToWorkspaceProcedure,
			connect.WithSchema(tenantServiceMethods.ByName("AssignUserToWorkspace")),
			connect.WithClientOptions(opts...),
		),
	}
}

// tenantServiceClient implements TenantServiceClient.
type tenantServiceClient struct {
	createWorkspace            *connect.Client[api.CreateWorkspaceRequest, api.CreateWorkspaceResponse]
	getWorkspaces              *connect.Client[api.GetWorkspacesRequest, api.GetWorkspacesResponse]
	editWorkspace              *connect.Client[api.EditWorkspaceRequest, api.EditWorkspaceResponse]
	getUsers                   *connect.Client[api.GetUsersRequest, api.GetUsersResponse]
	getUserWorkspaceAssignment *connect.Client[api.GetUserWorkspaceAssignmentRequest, api.GetUserWorkspaceAssignmentResponse]
	assignUserToWorkspace      *connect.Client[api.AssignUserToWorkspaceRequest, api.AssignUserToWorkspaceResponse]
}

// CreateWorkspace calls tenant.TenantService.CreateWorkspace.
func (c *tenantServiceClient) CreateWorkspace(ctx context.Context, req *connect.Request[api.CreateWorkspaceRequest]) (*connect.Response[api.CreateWorkspaceResponse], error) {
	return c.createWorkspace.CallUnary(ctx, req)
}

// GetWorkspaces calls tenant.TenantService.GetWorkspaces.
func (c *tenantServiceClient) GetWorkspaces(ctx context.Context, req *connect.Request[api.GetWorkspacesRequest]) (*connect.Response[api.GetWorkspacesResponse], error) {
	return c.getWorkspaces.CallUnary(ctx, req)
}

// EditWorkspace calls tenant.TenantService.EditWorkspace.
func (c *tenantServiceClient) EditWorkspace(ctx context.Context, req *connect.Request[api.EditWorkspaceRequest]) (*connect.Response[api.EditWorkspaceResponse], error) {
	return c.editWorkspace.CallUnary(ctx, req)
}

// GetUsers calls tenant.TenantService.GetUsers.
func (c *tenantServiceClient) GetUsers(ctx context.Context, req *connect.Request[api.GetUsersRequest]) (*connect.Response[api.GetUsersResponse], error) {
	return c.getUsers.CallUnary(ctx, req)
}

// GetUserWorkspaceAssignment calls tenant.TenantService.GetUserWorkspaceAssignment.
func (c *tenantServiceClient) GetUserWorkspaceAssignment(ctx context.Context, req *connect.Request[api.GetUserWorkspaceAssignmentRequest]) (*connect.Response[api.GetUserWorkspaceAssignmentResponse], error) {
	return c.getUserWorkspaceAssignment.CallUnary(ctx, req)
}

// AssignUserToWorkspace calls tenant.TenantService.AssignUserToWorkspace.
func (c *tenantServiceClient) AssignUserToWorkspace(ctx context.Context, req *connect.Request[api.AssignUserToWorkspaceRequest]) (*connect.Response[api.AssignUserToWorkspaceResponse], error) {
	return c.assignUserToWorkspace.CallUnary(ctx, req)
}

// TenantServiceHandler is an implementation of the tenant.TenantService service.
type TenantServiceHandler interface {
	CreateWorkspace(context.Context, *connect.Request[api.CreateWorkspaceRequest]) (*connect.Response[api.CreateWorkspaceResponse], error)
	GetWorkspaces(context.Context, *connect.Request[api.GetWorkspacesRequest]) (*connect.Response[api.GetWorkspacesResponse], error)
	EditWorkspace(context.Context, *connect.Request[api.EditWorkspaceRequest]) (*connect.Response[api.EditWorkspaceResponse], error)
	GetUsers(context.Context, *connect.Request[api.GetUsersRequest]) (*connect.Response[api.GetUsersResponse], error)
	GetUserWorkspaceAssignment(context.Context, *connect.Request[api.GetUserWorkspaceAssignmentRequest]) (*connect.Response[api.GetUserWorkspaceAssignmentResponse], error)
	AssignUserToWorkspace(context.Context, *connect.Request[api.AssignUserToWorkspaceRequest]) (*connect.Response[api.AssignUserToWorkspaceResponse], error)
}

// NewTenantServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewTenantServiceHandler(svc TenantServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	tenantServiceMethods := api.File_api_tenant_proto.Services().ByName("TenantService").Methods()
	tenantServiceCreateWorkspaceHandler := connect.NewUnaryHandler(
		TenantServiceCreateWorkspaceProcedure,
		svc.CreateWorkspace,
		connect.WithSchema(tenantServiceMethods.ByName("CreateWorkspace")),
		connect.WithHandlerOptions(opts...),
	)
	tenantServiceGetWorkspacesHandler := connect.NewUnaryHandler(
		TenantServiceGetWorkspacesProcedure,
		svc.GetWorkspaces,
		connect.WithSchema(tenantServiceMethods.ByName("GetWorkspaces")),
		connect.WithHandlerOptions(opts...),
	)
	tenantServiceEditWorkspaceHandler := connect.NewUnaryHandler(
		TenantServiceEditWorkspaceProcedure,
		svc.EditWorkspace,
		connect.WithSchema(tenantServiceMethods.ByName("EditWorkspace")),
		connect.WithHandlerOptions(opts...),
	)
	tenantServiceGetUsersHandler := connect.NewUnaryHandler(
		TenantServiceGetUsersProcedure,
		svc.GetUsers,
		connect.WithSchema(tenantServiceMethods.ByName("GetUsers")),
		connect.WithHandlerOptions(opts...),
	)
	tenantServiceGetUserWorkspaceAssignmentHandler := connect.NewUnaryHandler(
		TenantServiceGetUserWorkspaceAssignmentProcedure,
		svc.GetUserWorkspaceAssignment,
		connect.WithSchema(tenantServiceMethods.ByName("GetUserWorkspaceAssignment")),
		connect.WithHandlerOptions(opts...),
	)
	tenantServiceAssignUserToWorkspaceHandler := connect.NewUnaryHandler(
		TenantServiceAssignUserToWorkspaceProcedure,
		svc.AssignUserToWorkspace,
		connect.WithSchema(tenantServiceMethods.ByName("AssignUserToWorkspace")),
		connect.WithHandlerOptions(opts...),
	)
	return "/tenant.TenantService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case TenantServiceCreateWorkspaceProcedure:
			tenantServiceCreateWorkspaceHandler.ServeHTTP(w, r)
		case TenantServiceGetWorkspacesProcedure:
			tenantServiceGetWorkspacesHandler.ServeHTTP(w, r)
		case TenantServiceEditWorkspaceProcedure:
			tenantServiceEditWorkspaceHandler.ServeHTTP(w, r)
		case TenantServiceGetUsersProcedure:
			tenantServiceGetUsersHandler.ServeHTTP(w, r)
		case TenantServiceGetUserWorkspaceAssignmentProcedure:
			tenantServiceGetUserWorkspaceAssignmentHandler.ServeHTTP(w, r)
		case TenantServiceAssignUserToWorkspaceProcedure:
			tenantServiceAssignUserToWorkspaceHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedTenantServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedTenantServiceHandler struct{}

func (UnimplementedTenantServiceHandler) CreateWorkspace(context.Context, *connect.Request[api.CreateWorkspaceRequest]) (*connect.Response[api.CreateWorkspaceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("tenant.TenantService.CreateWorkspace is not implemented"))
}

func (UnimplementedTenantServiceHandler) GetWorkspaces(context.Context, *connect.Request[api.GetWorkspacesRequest]) (*connect.Response[api.GetWorkspacesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("tenant.TenantService.GetWorkspaces is not implemented"))
}

func (UnimplementedTenantServiceHandler) EditWorkspace(context.Context, *connect.Request[api.EditWorkspaceRequest]) (*connect.Response[api.EditWorkspaceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("tenant.TenantService.EditWorkspace is not implemented"))
}

func (UnimplementedTenantServiceHandler) GetUsers(context.Context, *connect.Request[api.GetUsersRequest]) (*connect.Response[api.GetUsersResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("tenant.TenantService.GetUsers is not implemented"))
}

func (UnimplementedTenantServiceHandler) GetUserWorkspaceAssignment(context.Context, *connect.Request[api.GetUserWorkspaceAssignmentRequest]) (*connect.Response[api.GetUserWorkspaceAssignmentResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("tenant.TenantService.GetUserWorkspaceAssignment is not implemented"))
}

func (UnimplementedTenantServiceHandler) AssignUserToWorkspace(context.Context, *connect.Request[api.AssignUserToWorkspaceRequest]) (*connect.Response[api.AssignUserToWorkspaceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("tenant.TenantService.AssignUserToWorkspace is not implemented"))
}
