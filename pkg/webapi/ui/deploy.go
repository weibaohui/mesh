package ui

import (
	"github.com/emicklei/go-restful"
	"github.com/sirupsen/logrus"
	"github.com/weibaohui/mesh/pkg/constants"
	"github.com/weibaohui/mesh/pkg/server"
	"github.com/weibaohui/mesh/pkg/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Inject(request *restful.Request, response *restful.Response) {
	ns := request.PathParameter("ns")
	name := request.PathParameter("name")
	mCtx := server.GlobalContext()
	deployment, err := mCtx.Apps.Apps().V1().Deployment().Get(ns, name, v1.GetOptions{})
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
