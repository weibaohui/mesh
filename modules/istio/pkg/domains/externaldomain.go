package domains

import (
	"fmt"
	"github.com/weibaohui/mesh/pkg/constants"
)

func GetPublicGateway(systemNamespace string) string {
	return fmt.Sprintf("%s.%s.svc.cluster.local", constants.MeshGateway, systemNamespace)
}

func GetExternalDomain(name, namespace, clusterDomain string) string {
	return fmt.Sprintf("%s-%s.%s", name, namespace, clusterDomain)
}
