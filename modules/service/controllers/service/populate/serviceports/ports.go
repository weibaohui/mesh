package serviceports

import (
	"fmt"
	"strings"

	meshv1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func Protocol(proto meshv1.Protocol) (protocol v1.Protocol) {
	switch proto {
	case meshv1.ProtocolUDP:
		protocol = v1.ProtocolUDP
	case meshv1.ProtocolSCTP:
		protocol = v1.ProtocolSCTP
	default:
		protocol = v1.ProtocolTCP
	}

	return
}

func ServiceNamedPorts(service *meshv1.Service) []v1.ServicePort {
	var (
		servicePorts []v1.ServicePort
	)

	ports := service.Spec.Ports
	for _, container := range service.Spec.Sidecars {
		ports = append(ports, container.Ports...)
	}

	portMap := map[string]meshv1.ContainerPort{}
	for _, port := range ports {
		portMap[fmt.Sprintf("%v/%v", port.Port, port.Protocol)] = port
	}

	for _, port := range portMap {
		if port.Port == 0 {
			port.Port = port.TargetPort
		}
		servicePort := v1.ServicePort{
			Name:     port.Name,
			Port:     port.Port,
			Protocol: Protocol(port.Protocol),
			TargetPort: intstr.IntOrString{
				IntVal: port.TargetPort,
			},
		}

		if servicePort.Name == "" {
			if port.Protocol == "" {
				port.Protocol = meshv1.ProtocolHTTP
			}
			servicePort.Name = strings.ToLower(fmt.Sprintf("%s-%d", port.Protocol, port.Port))
		}

		servicePorts = append(servicePorts, servicePort)
	}

	return servicePorts
}
