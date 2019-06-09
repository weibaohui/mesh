package populate

import "github.com/knative/pkg/apis/istio/v1alpha3"

type Protocol string

const (
	ProtocolTCP   Protocol = "TCP"
	ProtocolUDP   Protocol = "UDP"
	ProtocolSCTP  Protocol = "SCTP"
	ProtocolHTTP  Protocol = "HTTP"
	ProtocolHTTP2 Protocol = "HTTP2"
	ProtocolGRPC  Protocol = "GRPC"
)


var (
	supportedProtocol = []v1alpha3.PortProtocol{
		v1alpha3.ProtocolHTTP,
		// v1alpha3.ProtocolTCP,
		// v1alpha3.ProtocolGRPC,
		// v1alpha3.ProtocolHTTP2,
	}
)
