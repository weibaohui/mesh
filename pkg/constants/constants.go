package constants

var (
	DefaultHTTPOpenPort      = "80"
	DefaultHTTPSOpenPort     = "443"
	UseHostPort              = false
	UseIPAddress             = ""
	ServiceCidr              = ""
	DefaultServiceVersion    = "v0"
	IstioVersion             = "1.1.7"
	IstioGateway             = "ingressgateway"
	IstioMeshConfigKey       = "meshConfig"
	IstionConfigMapName      = "mesh"
	IstioSidecarTemplateName = "sidecarTemplate"
	Prometheus               = "prometheus"
	MeshGateway              = "mesh"
	ClusterDomainName        = "cluster.local"
	IstioInjector            = "istio-injecter"
	IstioInjectionEnable     = "istioInjectionEnable"
	IstioProxy               = "istio-proxy"
	DefaultShellExecCommand  = []string{"/bin/sh", "-c", `TERM=xterm-256color; export TERM; [ -x /bin/bash ] && ([ -x /usr/bin/script ] && /usr/bin/script -q -c "/bin/bash" /dev/null || exec /bin/bash) || exec /bin/sh`}
)
