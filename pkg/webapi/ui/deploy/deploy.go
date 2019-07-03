package deploy

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/sirupsen/logrus"
	"github.com/weibaohui/mesh/pkg/constants"
	"github.com/weibaohui/mesh/pkg/server"
	"github.com/weibaohui/mesh/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sort"
)

type Info struct {
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Ready      string `json:"ready"`
	MeshEnable string `json:"meshEnable"`
}

func List(request *restful.Request, response *restful.Response) {
	ns := request.QueryParameter("ns")
	mCtx := server.GlobalContext()
	list, err := mCtx.Apps.Apps().V1().Deployment().Cache().List(ns, labels.Everything())
	if err != nil {
		logrus.Errorf("List Deployments %s:%s", ns, err.Error())
		response.WriteError(500, err)
		return
	}

	var infos []Info
	for _, p := range list {
		infos = append(infos, Info{
			Name:       p.Name,
			Namespace:  p.Namespace,
			Ready:      fmt.Sprintf("%d/%d", p.Status.AvailableReplicas, p.Status.Replicas),
			MeshEnable: utils.GetValueFrom(p.GetAnnotations(), constants.IstioInjectionEnable),
		})
	}
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Name < infos[j].Name
	})
	i := struct {
		Code  int    `json:"code"`
		Count int    `json:"count"`
		Data  []Info `json:"data"`
		Msg   string `json:"msg"`
	}{
		Code:  0,
		Count: len(infos),
		Data:  infos,
		Msg:   "",
	}

	response.WriteAsJson(i)
}

func Inject(request *restful.Request, response *restful.Response) {
	ns := request.PathParameter("ns")
	name := request.PathParameter("name")
	mCtx := server.GlobalContext()
	deployment, err := mCtx.Apps.Apps().V1().Deployment().Get(ns, name, metav1.GetOptions{})
	if err != nil {
		logrus.Errorf("%s-%s:%s", ns, name, err.Error())
		response.WriteError(500, err)
		return
	}
	annotations := deployment.GetAnnotations()
	annotations = utils.Merge(annotations, map[string]string{
		constants.IstioInjectionEnable: "true",
	})
	deployment.SetAnnotations(annotations)
	update, err := mCtx.Apps.Apps().V1().Deployment().Update(deployment)

	if err != nil {
		logrus.Errorf("%s-%s:%s", ns, name, err.Error())
		response.WriteError(500, err)
		return
	}
	response.WriteEntity(update)
}
