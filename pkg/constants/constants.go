package constants

var (
	DefaultHTTPOpenPort      = "80"
	DefaultHTTPSOpenPort     = "443"
	UseHostPort              = false
	UseIPAddress             = ""
	ServiceCidr              = ""
	DefaultServiceVersion    = "v0"
	GatewaySecretName        = "mesh-certs"
	IstioVersion             = "1.1.7"
	IstioGateway             = "ingressgateway"
	IstioMeshConfigKey       = "meshConfig"
	IstionConfigMapName      = "mesh"
	IstioSidecarTemplateName = "sidecarTemplate"
	IstioStackName           = "istio"
	IstioTelemetry           = "istio-telemetry"
	ProductionType           = "production"
	Prometheus               = "prometheus"
	MeshGateway              = "mesh-gateway"
	StagingType              = "staging"
	ClusterDomainName        = "cluster.local"
)
