package constants

var (
	DefaultHTTPOpenPort      = "9080"
	DefaultHTTPSOpenPort     = "9443"
	UseHostPort              = false
	UseIPAddress             = ""
	ServiceCidr              = ""
	DefaultServiceVersion    = "v0"
	GatewaySecretName        = "mesh-certs"
	IstioGateway             = "istio-gateway"
	IstioMeshConfigKey       = "meshConfig"
	IstionConfigMapName      = "mesh"
	IstioSidecarTemplateName = "sidecarTemplate"
	IstioStackName           = "istio"
	IstioTelemetry           = "istio-telemetry"
	ProductionType           = "production"
	Prometheus               = "prometheus"
	MeshGateway              = "mesh-gateway"
	StagingType              = "staging"
)
