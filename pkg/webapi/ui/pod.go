package ui

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/sirupsen/logrus"
	"github.com/weibaohui/mesh/pkg/server"
	"k8s.io/apimachinery/pkg/labels"
)

type simplePodInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Ready     string `json:"ready"`
	PodIP     string `json:"pod_ip"`
	Status    string `json:"status"`
	Restart   string `json:"restart"`
	Age       string `json:"age"`
}

func ListPod(request *restful.Request, response *restful.Response) {
	ns := request.QueryParameter("ns")
	logrus.Println(ns)
	mCtx := server.GlobalContext()
	list, err := mCtx.Core.Core().V1().Pod().Cache().List(ns, labels.Everything())
	if err != nil {
		fmt.Println(err.Error())
	}
	var podlist []simplePodInfo
	for _, p := range list {
		podlist = append(podlist, simplePodInfo{
			Name:      p.Name,
			Namespace: p.Namespace,
			Ready:     p.Status.Reason,
			Status:    string(p.Status.Phase),
			PodIP:     p.Status.PodIP,
			Restart:   "",
			Age:       p.CreationTimestamp.String(),
		})
	}
	i := struct {
		Code  int             `json:"code"`
		Count int             `json:"count"`
		Data  []simplePodInfo `json:"data"`
		Msg   string          `json:"msg"`
	}{
		Code:  0,
		Count: len(podlist),
		Data:  podlist,
		Msg:   "",
	}

	response.WriteAsJson(i)
}
