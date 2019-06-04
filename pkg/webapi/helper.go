package webapi

import (
	"context"
	"flag"
	"fmt"
	"github.com/weibaohui/mesh/pkg/generated/controllers/mesh.oauthd.com"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typeV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"sync"
)

var cli kubernetes.Interface
var meshCli mesh.Interface
var once = sync.Once{}

type Helper struct {
	cli     kubernetes.Interface
	meshCli mesh.Interface
}

func NewHelper() *Helper {
	once.Do(func() {
		meshCli = GetMeshClient()
	})

	return &Helper{meshCli: meshCli}
}

func (h *Helper) MeshClient() mesh.Interface {

	return h.meshCli
}
func (h *Helper) Pods(ns string) typeV1.PodInterface {
	return h.cli.CoreV1().Pods(ns)
}
func (h *Helper) Services(ns string) typeV1.ServiceInterface {
	return h.cli.CoreV1().Services(ns)
}

func (h *Helper) GetPod(ns, podName string) (*coreV1.Pod, error) {
	return h.Pods(ns).Get(podName, metaV1.GetOptions{})
}

func (h *Helper) GetService(ns, svcName string) (*coreV1.Service, error) {
	return h.Services(ns).Get(svcName, metaV1.GetOptions{})
}
func (h *Helper) IsServiceExists(ns, svcName string) bool {
	_, e := h.Services(ns).Get(svcName, metaV1.GetOptions{})
	if e == nil {
		return true
	}
	return false
}

func GetClient() kubernetes.Interface {
	var kubeConfig *string
	if home := homeDir(); home != "" {
		s := filepath.Join(home, ".kube", "config")
		kubeConfig = flag.String("kubeconfig", s, "kubeconfig存放位置")
	} else {
		kubeConfig = flag.String("kubeconfig", "", "kubeconfig存放位置")
	}
	flag.Parse()
	var config *rest.Config
	var err error
	config, err = rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeConfig)
	}

	if err != nil {
		panic(err.Error())
	}

	cli, e := kubernetes.NewForConfig(config)
	if e != nil {
		panic(e.Error())
	}
	return cli

}

func GetMeshClient() mesh.Interface {
	var kubeConfig = "/Users/weibh/.kube/config"

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)

	if err != nil {
		fmt.Println(err)
	}
	factory := mesh.NewFactoryFromConfigOrDie(config)
	background := context.Background()
	err = factory.Start(background, 5)

	if err != nil {
		fmt.Println(err)
	}

	return factory.Mesh()

}

func homeDir() string {
	if s := os.Getenv("HOME"); s != "" {
		return s
	}
	return os.Getenv("USERPROFILE")
}
