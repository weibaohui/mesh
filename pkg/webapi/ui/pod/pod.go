package pod

import (
	"fmt"
	"sort"
	"strings"

	"github.com/emicklei/go-restful"
	"github.com/weibaohui/mesh/pkg/constants"
	"github.com/weibaohui/mesh/pkg/server"
	"github.com/weibaohui/mesh/pkg/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

// Info pod info
type Info struct {
	Deploy     string `json:"deploy"`
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Ready      string `json:"ready"`
	PodIP      string `json:"podIp"`
	Status     string `json:"status"`
	Restart    int32  `json:"restart"`
	MeshEnable string `json:"meshEnable"`
	Containers string `json:"containers"`
}

func buildStatus(pod *v1.Pod) string {
	for _, cs := range pod.Status.ContainerStatuses {

		if cs.State.Waiting != nil {
			return cs.State.Waiting.Reason
		}
		if cs.State.Terminated != nil {
			return cs.State.Terminated.Reason
		}
		if cs.State.Running != nil {
			return string(pod.Status.Phase)
		}
	}
	// for _, c := range pod.Status.ContainerStatuses {
	// 	if c.State.Waiting != nil {
	// 		return c.State.Waiting.Reason
	// 	}
	// 	if c.State.Terminated != nil {
	// 		return c.State.Terminated.Reason
	// 	}
	// }
	return ""
}

func buildReadyStatusCount(pod *v1.Pod) string {
	var c int32
	for _, cs := range pod.Status.ContainerStatuses {
		if cs.Ready {
			c++
		}

	}
	return fmt.Sprintf("%d/%d", c, len(pod.Spec.Containers))
}
func buildRestartCount(pod *v1.Pod) int32 {
	var c int32
	for _, cs := range pod.Status.ContainerStatuses {
		c += cs.RestartCount
	}
	return c
}

func buildContainers(d *v1.Pod) string {
	var cNames string
	for _, c := range d.Spec.Containers {
		cNames = cNames + c.Name + ","
	}
	cNames = strings.TrimSuffix(cNames, ",")
	return cNames
}

// List 列表显示pod 信息
func List(request *restful.Request, response *restful.Response) {
	ns := request.QueryParameter("ns")
	appName := request.QueryParameter("appName")
	selector := labels.Everything()
	if len(appName) > 0 {
		r, _ := labels.NewRequirement("app", selection.Equals, []string{appName})
		selector = labels.NewSelector().Add(*r)
	}

	mCtx := server.GlobalContext()
	list, err := mCtx.Core.Core().V1().Pod().Cache().List(ns, selector)
	if err != nil {
		fmt.Println(err.Error())
	}
	var infos []Info
	for _, p := range list {
		infos = append(infos, Info{
			Deploy:     utils.GetValueFrom(p.GetLabels(), "app"),
			Name:       p.Name,
			Namespace:  p.Namespace,
			Ready:      buildReadyStatusCount(p),
			PodIP:      p.Status.PodIP,
			Status:     buildStatus(p),
			Restart:    buildRestartCount(p),
			MeshEnable: utils.GetValueFrom(p.GetAnnotations(), constants.IstioInjectionEnable),
			Containers: buildContainers(p),
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
