package ui

import (
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/gorilla/websocket"
	"github.com/weibaohui/mesh/pkg/server"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"net/http"
)

type terminal struct {
	conn          *websocket.Conn
	size          chan *remotecommand.TerminalSize
	ns            string
	podName       string
	containerName string
}

func (t *terminal) Read(p []byte) (n int, err error) {
	_, ps, err := t.conn.ReadMessage()
	if err != nil {
		return 0, err
	}
	size := remotecommand.TerminalSize{}
	if err = json.Unmarshal(ps, &size); err == nil {
		t.size <- &size
		return 0, nil
	} else {
		fmt.Println(err.Error())
		return copy(p, ps), nil
	}
}
func (t *terminal) Write(p []byte) (n int, err error) {
	writer, err := t.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return 0, nil
	}
	defer writer.Close()
	return writer.Write(p)
}
func (t *terminal) Next() *remotecommand.TerminalSize {
	size := <-t.size
	fmt.Println("size", size)
	return size
}

// 	ws.Route(ws.GET("/ns/{ns}/podName/{podName}/exec")
func PodExec(request *restful.Request, response *restful.Response) {
	params := request.PathParameters()
	ns := params["ns"]
	podName := params["podName"]
	containerName := request.QueryParameter("containerName")
	if containerName == "" {
		// 没有指定，获取第一个
		containerName, _ = GetFirstContainerName(ns, podName)
	}

	c, err := upgrader.Upgrade(response, request.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	t := &terminal{
		conn:          c,
		ns:            ns,
		podName:       podName,
		containerName: containerName,
	}

	err = executor(t)
	if err != nil {
		fmt.Println(err.Error())
	}

}

func executor(t *terminal) error {
	var defaultCommand = []string{"/bin/sh", "-c", `TERM=xterm-256color; export TERM; [ -x /bin/bash ] && ([ -x /usr/bin/script ] && /usr/bin/script -q -c "/bin/bash" /dev/null || exec /bin/bash) || exec /bin/sh`};
	mctx := server.GlobalContext()
	config := server.GlobalKubeConfig()
	req := mctx.K8s.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(t.podName).
		Namespace(t.ns).
		SubResource("exec").
		VersionedParams(
			&v1.PodExecOptions{
				TypeMeta:  v12.TypeMeta{},
				Stdin:     true,
				Stdout:    true,
				TTY:       true,
				Container: t.containerName,
				Command:   defaultCommand,
			},
			scheme.ParameterCodec,
		)
	restConfig, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		fmt.Errorf("clientcmd.BuildConfigFromFlags =%s ", err.Error())
		return err
	}
	exec, err := remotecommand.NewSPDYExecutor(restConfig, http.MethodPost, req.URL())
	if err != nil {
		fmt.Errorf("remotecommand.NewSPDYExecutor =%s ", err.Error())
		return err
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:             t,
		Stdout:            t,
		Stderr:            t,
		Tty:               true,
		TerminalSizeQueue: t,
	})
	return err

}
