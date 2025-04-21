package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/MorhafAlshibly/iunvi/gen/api/apiconnect"
	"github.com/MorhafAlshibly/iunvi/internal/model"
	"github.com/MorhafAlshibly/iunvi/pkg/middleware"
	"github.com/MorhafAlshibly/iunvi/pkg/sculpt"
	_ "github.com/microsoft/go-mssqldb"
	"github.com/microsoft/go-mssqldb/azuread"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
)

var (
	// Flags set from command line/environment variables
	fs                     = ff.NewFlagSet("model")
	port                   = fs.Uint('p', "port", 8080, "the default port to listen on")
	origin                 = fs.StringLong("origin", "http://localhost:7575", "the default origin to allow")
	subscriptionId         = fs.StringLong("subscriptionId", "", "Azure subscription ID")
	resourceGroupName      = fs.StringLong("resourceGroupName", "rg-iunvi-dev-eastus-001", "Azure resource group name")
	tenantId               = fs.StringLong("tenantId", "", "Azure AD tenant ID")
	clientId               = fs.StringLong("clientId", "", "Azure AD client ID")
	audience               = fs.StringLong("audience", "", "Azure AD audience")
	clientSecret           = fs.StringLong("clientSecret", "", "Azure AD client secret")
	jwks                   = fs.StringLong("jwks", "https://login.microsoftonline.com/common/discovery/v2.0/keys", "Azure AD JWKS URL")
	sqlServer              = fs.StringLong("sqlServer", "", "SQL Server")
	sqlDatabase            = fs.StringLong("sqlDatabase", "", "SQL Database")
	registryName           = fs.StringLong("registryName", "", "Azure Container Registry URL")
	registryTokenPrefix    = fs.StringLong("registryTokenPrefix", "webapp", "Azure Container Registry Token Prefix")
	storageAccountName     = fs.StringLong("storageAccountName", "", "Azure Storage Account Name")
	modelRunsContainerName = fs.StringLong("modelRunsContainerName", "model-runs", "Azure Storage Model Runs Container Name")
	kubeConfigPath         = fs.StringLong("kubeConfigPath", "", "Kube Config Path")
	inputFileMountPath     = fs.StringLong("inputFileMountPath", "/mnt/input", "Input file mount path")
	parametersMountPath    = fs.StringLong("parametersMountPath", "/mnt/parameters", "Parameters mount path")
	outputFileMountPath    = fs.StringLong("outputFileMountPath", "/mnt/output", "Output file mount path")
	fileGroupsPVCName      = fs.StringLong("fileGroupsPVCName", "", "File groups PVC name")
	modelRunsPVCName       = fs.StringLong("modelRunsPVCName", "", "Model runs PVC name")
)

func main() {
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("MODEL"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))
	if err != nil {
		fmt.Printf("%s\n", ffhelp.Flags(fs))
		fmt.Printf("failed to parse flags: %v", err)
		return
	}
	db, err := sql.Open(azuread.DriverName, sculpt.DatabaseUrl(*sqlServer, *sqlDatabase, *clientId, *clientSecret))
	if err != nil {
		fmt.Printf("failed to open database connection: %v", err)
		return
	}
	defer db.Close()
	service := model.NewService(
		model.WithSubscriptionId(*subscriptionId),
		model.WithResourceGroupName(*resourceGroupName),
		model.WithTenantId(*tenantId),
		model.WithClientId(*clientId),
		model.WithClientSecret(*clientSecret),
		model.WithRegistryName(*registryName),
		model.WithRegistryTokenPrefix(*registryTokenPrefix),
		model.WithStorageAccountName(*storageAccountName),
		model.WithModelRunsContainerName(*modelRunsContainerName),
		model.WithKubeConfigPath(*kubeConfigPath),
		model.WithInputFileMountPath(*inputFileMountPath),
		model.WithParametersMountPath(*parametersMountPath),
		model.WithOutputFileMountPath(*outputFileMountPath),
		model.WithFileGroupsPVCName(*fileGroupsPVCName),
		model.WithModelRunsPVCName(*modelRunsPVCName),
	)
	mux := http.NewServeMux()
	path, handler := apiconnect.NewModelServiceHandler(service)
	cors := middleware.NewCORS(middleware.WithAllowedOrigins([]string{*origin}))
	auth := middleware.NewAuthentication(
		middleware.WithAudience(*audience),
		middleware.WithJWKS(*jwks),
	)
	transaction := middleware.NewTransaction(middleware.WithDB(db))
	mux.Handle(path, cors.Middleware(transaction.Middleware(auth.Middleware(handler))))
	if err := http.ListenAndServe(
		fmt.Sprintf(":%d", *port),
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		fmt.Printf("failed to start server: %v", err)
		return
	}
}
