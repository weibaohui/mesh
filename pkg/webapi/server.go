package webapi

import (
	"fmt"
	"github.com/emicklei/go-restful"
	v1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
	"io/ioutil"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"log"
	"net/http"
)

func Start() {
	container := restful.NewContainer()
	ws := new(restful.WebService)
	ws.Route(ws.POST("/version").
		To(ports).
		Produces(restful.MIME_JSON))
	container.Add(ws)
	log.Fatal(http.ListenAndServe(":9999", container))

}

// GET /ports
func ports(req *restful.Request, resp *restful.Response) {
	bytes, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(bytes))
	app := v1.App{

		Spec: v1.AppSpec{

			Revisions: []v1.Revision{
				{
					Public:      true,
					ServiceName: "whoami-v2",
					Version:     "v2",
					Weight:      50,
				},
				{
					Public:      true,
					ServiceName: "whoami-v3",
					Version:     "v3",
					Weight:      50,
				},
			},
		},
	}

	controller := NewHelper().MeshClient().V1().App()
	obj, err := controller.Get("default", "whoami", v12.GetOptions{})
	if err != nil {
		fmt.Println(err.Error())
	}
	obj.Spec=app.Spec

	appRet, err := controller.Update(obj)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(appRet.CreationTimestamp)

	resp.WriteEntity("ok")
}
