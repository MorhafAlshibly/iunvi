package sculpt

import (
	"fmt"
	"strings"
)

func RegistryScope(workspaceId string) string {
	return fmt.Sprintf("scope-%s", strings.ToLower(workspaceId))
}

func RegistryTokenName(workspaceId string, registryTokenPrefix string) string {
	return fmt.Sprintf("%s-%s", registryTokenPrefix, strings.ToLower(workspaceId))
}
