package sculpt

import (
	"fmt"
)

func DatabaseUrl(sqlServer string, sqlDatabase string, clientId string, clientSecret string) string {
	return fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;encrypt=True;fedauth=ActiveDirectoryServicePrincipal;", sqlServer, clientId, clientSecret, 1433, sqlDatabase)
}
