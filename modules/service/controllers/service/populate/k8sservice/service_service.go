package k8sservice

import (
	"strconv"

	"github.com/rancher/wrangler/pkg/objectset"
	"github.com/weibaohui/mesh/modules/service/controllers/service/populate/servicelabels"
	"github.com/weibaohui/mesh/modules/service/controllers/service/populate/serviceports"
	meshv1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	"github.com/weibaohui/mesh/pkg/constants"
	"github.com/weibaohui/mesh/pkg/constructors"
	"github.com/weibaohui/mesh/pkg/services"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func serviceSelector(service *meshv1.Service, systemNamespace string, os *objectset.ObjectSet) {
	labels := servicelabels.ServiceLabels(service)
	selectorLabels := servicelabels.SelectorLabels(service)
	app, version := services.AppAndVersion(service)
	svc := newServiceSelector(app+"-"+version, service.Namespace, v1.ServiceTypeClusterIP, labels, selectorLabels)
	if isGateway(service, systemNamespace) && !constants.UseHostPort {
		svc.Spec.Type = v1.ServiceTypeLoadBalancer
		httpPort, _ := strconv.Atoi(constants.DefaultHTTPOpenPort)
		httpsPort, _ := strconv.Atoi(constants.DefaultHTTPSOpenPort)
		svc.Spec.Ports = []v1.ServicePort{
			{
				Name:       "http",
				Protocol:   v1.ProtocolTCP,
				TargetPort: intstr.FromInt(httpPort),
				Port:       int32(httpPort),
			},
			{
				Name:       "https",
				Protocol:   v1.ProtocolTCP,
				TargetPort: intstr.FromInt(httpsPort),
				Port:       int32(httpsPort),
			},
		}
	} else {
		if len(serviceports.ServiceNamedPorts(service)) > 0 {
			svc.Spec.Ports = serviceports.ServiceNamedPorts(service)
		}
	}
	os.Add(svc)
}

func isGateway(service *meshv1.Service, systemNamespace string) bool {
	return service.Name == constants.IstioGateway && service.Namespace == systemNamespace
}

func newServiceSelector(name, namespace string, serviceType v1.ServiceType, labels, selectorLabels map[string]string) *v1.Service {
	return constructors.NewService(namespace, name, v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Labels: labels,
		},
		Spec: v1.ServiceSpec{
			Type:     serviceType,
			Selector: selectorLabels,
			Ports: []v1.ServicePort{
				{
					Name:       "default",
					Protocol:   v1.ProtocolTCP,
					TargetPort: intstr.FromInt(80),
					Port:       80,
				},
			},
		},
	})
}
