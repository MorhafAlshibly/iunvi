package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/MorhafAlshibly/iunvi/gen/api/apiconnect"
	"github.com/MorhafAlshibly/iunvi/internal/tenantManagement"
	"github.com/MorhafAlshibly/iunvi/pkg/middleware"
	_ "github.com/microsoft/go-mssqldb"
	"github.com/microsoft/go-mssqldb/azuread"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
)

var (
	// Flags set from command line/environment variables
	fs                     = ff.NewFlagSet("tenantManagement")
	port                   = fs.Uint('p', "port", 8080, "the default port to listen on")
	azureSubscriptionId    = fs.StringLong("azureSubscriptionId", "", "Azure subscription ID")
	azureResourceGroupName = fs.StringLong("azureResourceGroupName", "rg-iunvi-dev-eastus-001", "Azure resource group name")
	azureAdTenantID        = fs.StringLong("azureAdTenantID", "", "Azure AD tenant ID")
	azureAdClientID        = fs.StringLong("azureAdClientID", "", "Azure AD client ID")
	azureAdAudience        = fs.StringLong("azureAdAudience", "", "Azure AD audience")
	azureAdClientSecret    = fs.StringLong("azureAdClientSecret", "", "Azure AD client secret")
	azureAdJWKS            = fs.StringLong("azureAdJWKS", "https://login.microsoftonline.com/common/discovery/v2.0/keys", "Azure AD JWKS URL")
	sqlServer              = fs.StringLong("sqlServer", "", "SQL Server")
	sqlDatabase            = fs.StringLong("sqlDatabase", "", "SQL Database")
	registryName           = fs.StringLong("registryName", "criunvideveastus001", "Azure Container Registry URL")
	registryTokenName      = fs.StringLong("registryTokenName", "webapp", "Azure Container Registry Token Name")
	storageAccountName     = fs.StringLong("storageAccountName", "saiunvideveastus001", "Azure Storage Account Name")
	storageContainerName   = fs.StringLong("storageContainerName", "landing-zone", "Azure Storage Container Name")
)

func main() {
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("TENANTMANAGEMENT"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))
	if err != nil {
		fmt.Printf("%s\n", ffhelp.Flags(fs))
		fmt.Printf("failed to parse flags: %v", err)
		return
	}
	db, err := sql.Open(azuread.DriverName, fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;encrypt=True;fedauth=ActiveDirectoryServicePrincipal;", *sqlServer, *azureAdClientID, *azureAdClientSecret, 1433, *sqlDatabase))
	if err != nil {
		fmt.Printf("failed to open database connection: %v", err)
		return
	}
	defer db.Close()
	service := tenantManagement.NewService(
		tenantManagement.WithSubscriptionId(*azureSubscriptionId),
		tenantManagement.WithResourceGroupName(*azureResourceGroupName),
		tenantManagement.WithTenantId(*azureAdTenantID),
		tenantManagement.WithClientId(*azureAdClientID),
		tenantManagement.WithClientSecret(*azureAdClientSecret),
		tenantManagement.WithRegistryName(*registryName),
		tenantManagement.WithRegistryTokenName(*registryTokenName),
		tenantManagement.WithStorageAccountName(*storageAccountName),
		tenantManagement.WithStorageContainerName(*storageContainerName),
	)
	mux := http.NewServeMux()
	path, handler := apiconnect.NewTenantManagementServiceHandler(service)
	cors := middleware.NewCORS(middleware.WithAllowedOrigins([]string{"http://localhost:7575"}))
	auth := middleware.NewAuthentication(
		middleware.WithAudience(*azureAdAudience),
		middleware.WithAzureAdJWKS(*azureAdJWKS),
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
