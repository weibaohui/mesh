package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful"
	v1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	"github.com/weibaohui/mesh/pkg/server"
	"github.com/weibaohui/mesh/pkg/webapi/ui"
	"io/ioutil"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"time"

	"log"
	"net/http"
)

func Start(ctx context.Context) {
	container := restful.NewContainer()
	ws := new(restful.WebService)
	ws.Route(ws.POST("/version").To(ports).
		Produces(restful.MIME_JSON))
	ws.Route(ws.GET("/ns/{ns}/podName/{podName}/log").To(ui.GetContainerLog)).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/tty").To(ui.Tty))
	ws.Route(ws.GET("/tt").To(tt).Produces(restful.MIME_JSON))
	ws.Route(ws.GET("/pods").To(ui.ListPod).Produces(restful.MIME_JSON))
	ws.Route(ws.GET("/deploy/inject/{ns}/{name}").To(ui.Inject).Produces(restful.MIME_JSON))
	container.Add(ws)
	fmt.Println("SERVER 9999")

	// Add container filter to enable CORS
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST"},
		CookiesAllowed: false,
		Container:      container}
	container.Filter(cors.Filter)

	go sync(ctx)
	log.Fatal(http.ListenAndServe(":9999", container))
}

func sync(ctx context.Context) {

	ticker := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-ticker.C:
			// mCtx := server.GlobalContext()
			// mCtx.Core.Sync(ctx)
			// logrus.Info("sync")
		}
	}

}

func tt(request *restful.Request, response *restful.Response) {
	mCtx := server.GlobalContext()
	requirement, err := labels.NewRequirement("name", selection.NotIn, []string{"x"})
	selector := labels.NewSelector().Add(*requirement)
	list, err := mCtx.Apps.Apps().V1().Deployment().Cache().List("default", selector)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, d := range list {
		fmt.Println(d.Name)
	}
	response.WriteAsJson(list)
}

type instanceWeight struct {
	Instance string `json:"instance"`
	Weight   int    `json:"weight"`
	Version  string `json:"version"`
}

// GET /ports
func ports(req *restful.Request, resp *restful.Response) {
	bytes, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	var instances []instanceWeight

	// bytes, err = json.Marshal(instances)
	// fmt.Println(string(bytes))

	err = json.Unmarshal(bytes, &instances)
	if err != nil {
		fmt.Println(err.Error())
	}

	app := v1.App{

		Spec: v1.AppSpec{
			// Revisions: []v1.Revision{
			// 	{
			// 		Public:      true,
			// 		ServiceName: "whoami-v2",
			// 		Version:     "v2",
			// 		Weight:      50,
			// 	},
			// 	{
			// 		Public:      true,
			// 		ServiceName: "whoami-v3",
			// 		Version:     "v3",
			// 		Weight:      50,
			// 	},
			// },
		},
	}
	for _, v := range instances {
		revision := v1.Revision{
			Public:      true,
			ServiceName: v.Instance,
			Version:     v.Version,
			Weight:      v.Weight,
		}

		app.Spec.Revisions = append(app.Spec.Revisions, revision)

	}

	mCtx := server.GlobalContext()
	controller := mCtx.Mesh.Mesh().V1().App()
	obj, err := controller.Get("default", "whoami", v12.GetOptions{})
	if err != nil {
		fmt.Println(err.Error())
	}

	obj.Spec = app.Spec

	appRet, err := controller.Update(obj)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(len(appRet.Spec.Revisions))

	resp.WriteEntity("ok")
}
