# User manual

To deploy the platform, you first need to have Docker and Docker Compose installed on your machine. You can find the installation instructions for Docker [here](https://docs.docker.com/get-docker/) and for Docker Compose [here](https://docs.docker.com/compose/install/).

Before deploying the platform, make sure you have the following environment variables set in your `./env/.env.dev` file (fill in the values for your environment in Azure):

Note that the `./env/.env.dev` file is not included in the repository for security reasons. You will need to create this file yourself.

```env
# Environment variables
TENANT_PORT=8001
TENANT_ORIGIN=http://localhost:7575
TENANT_SUBSCRIPTIONID=
TENANT_RESOURCEGROUPNAME=
TENANT_TENANTID=
TENANT_CLIENTID=
TENANT_AUDIENCE=
TENANT_CLIENTSECRET=
TENANT_JWKS=https://login.microsoftonline.com/common/discovery/v2.0/keys
TENANT_SQLSERVER=
TENANT_SQLDATABASE=
TENANT_REGISTRYNAME=
TENANT_REGISTRYTOKENPREFIX=webapp

FILE_PORT=8002
FILE_ORIGIN=http://localhost:7575
FILE_TENANTID=
FILE_CLIENTID=
FILE_AUDIENCE=
FILE_CLIENTSECRET=
FILE_JWKS=https://login.microsoftonline.com/common/discovery/v2.0/keys
FILE_SQLSERVER=
FILE_SQLDATABASE=
FILE_STORAGEACCOUNTNAME=
FILE_LANDINGZONECONTAINERNAME=landing-zone
FILE_FILEGROUPSCONTAINERNAME=file-groups

MODEL_PORT=8003
MODEL_ORIGIN=http://localhost:7575
MODEL_SUBSCRIPTIONID=
MODEL_RESOURCEGROUPNAME=
MODEL_TENANTID=
MODEL_CLIENTID=
MODEL_AUDIENCE=
MODEL_CLIENTSECRET=
MODEL_JWKS=https://login.microsoftonline.com/common/discovery/v2.0/keys
MODEL_SQLSERVER=
MODEL_SQLDATABASE=
MODEL_REGISTRYNAME=
MODEL_REGISTRYTOKENPREFIX=webapp
MODEL_STORAGEACCOUNTNAME=
MODEL_MODELRUNSCONTAINERNAME=model-runs
MODEL_KUBECONFIGPATH=
MODEL_INPUTFILEMOUNTPATH=/mnt/input
MODEL_PARAMETERSMOUNTPATH=/mnt/parameters
MODEL_OUTPUTFILEMOUNTPATH=/mnt/output
MODEL_FILEGROUPSPVCNAME=pvc-file-groups-iunvi-dev-eastus-001
MODEL_MODELRUNSPVCNAME=pvc-model-runs-iunvi-dev-eastus-001

DASHBOARD_PORT=8004
DASHBOARD_ORIGIN=http://localhost:7575
DASHBOARD_TENANTID=
DASHBOARD_CLIENTID=
DASHBOARD_AUDIENCE=
DASHBOARD_CLIENTSECRET=
DASHBOARD_JWKS=https://login.microsoftonline.com/common/discovery/v2.0/keys
DASHBOARD_SQLSERVER=
DASHBOARD_SQLDATABASE=
DASHBOARD_REGISTRYNAME=
DASHBOARD_REGISTRYTOKENPREFIX=webapp
DASHBOARD_STORAGEACCOUNTNAME=
DASHBOARD_MODELRUNSCONTAINERNAME=model-runs
DASHBOARD_DASHBOARDSCONTAINERNAME=dashboards
DASHBOARD_MODELRUNDASHBOARDSCONTAINERNAME=model-run-dashboards
DASHBOARD_KUBECONFIGPATH=
DASHBOARD_APPLYDASHBOARDIMAGENAME=apply-dashboard
DASHBOARD_MODELRUNSPVCNAME=pvc-model-runs-iunvi-dev-eastus-001
DASHBOARD_DASHBOARDSPVCNAME=pvc-dashboards-iunvi-dev-eastus-001
DASHBOARD_MODELRUNDASHBOARDSPVCNAME=pvc-model-run-dashboards-iunvi-dev-eastus-001


VITE_CLIENTID=
VITE_REGISTRYNAME=
VITE_TENANTID=common
VITE_REDIRECTURI=http://localhost:7575
VITE_SCOPE=
VITE_TENANTURL=http://localhost:8001
VITE_FILEURL=http://localhost:8002
VITE_MODELURL=http://localhost:8003
VITE_DASHBOARDURL=http://localhost:8004
```

Once you have Docker and Docker Compose installed, you can deploy the platform by running the following command in the root directory of the project:

```bash
docker-compose up
```

This will start all the services defined in the `docker-compose.yml` file. The platform will be accessible at `http://localhost:7575`.
