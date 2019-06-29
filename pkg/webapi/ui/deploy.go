package ui

import (
	"github.com/emicklei/go-restful"
	"github.com/weibaohui/mesh/pkg/constants"
	"github.com/weibaohui/mesh/pkg/server"
	"github.com/weibaohui/mesh/pkg/utils"
)

func Inject(request *restful.Request, response *restful.Response) {
	ns := request.PathParameter("ns")
	name := request.PathParameter("name")
	mCtx := server.GlobalContext()
	deployment, err := mCtx.Apps.Apps().V1().Deployment().Cache().Get(ns, name)
	if err != nil {
		response.WriteError(500, err)
		return
	}
	annotations := deployment.GetAnnotations()
	annotations = utils.Merge(annotations, map[string]string{
		constants.IstioInjection: "true",
	})
	deployment.SetAnnotations(annotations)
	update, err := mCtx.Apps.Apps().V1().Deployment().Update(deployment)

	if err != nil {
		response.WriteError(500, err)
		return
	}
	response.WriteEntity(update)
}
