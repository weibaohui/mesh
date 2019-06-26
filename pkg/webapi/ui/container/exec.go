package container

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/gorilla/websocket"
	"github.com/weibaohui/mesh/pkg/constants"
	"github.com/weibaohui/mesh/pkg/server"
	"golang.org/x/sync/errgroup"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"net/http"
	"sync"
)

var terminalMaps sync.Map

func (t *terminal) saveTerminal() {
	terminalMaps.Store(t.key(), t)
}
func (t *terminal) removeTerminal() {
	terminalMaps.Delete(t.key())
}
func (t *terminal) getTerminal() (*terminal, bool) {
	if value, ok := terminalMaps.Load(t.key()); ok {
		return value.(*terminal), true
	}
	return nil, false
}
func (t *terminal) key() string {
	return fmt.Sprintf("%s/%s/%s", t.ns, t.podName, t.containerName)
}

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
	sizes := <-t.size
	return sizes
}

func Resize(req *restful.Request, resp *restful.Response) {
	t1 := &terminal{
		ns:            req.QueryParameter("ns"),
		podName:       req.QueryParameter("podName"),
		containerName: req.QueryParameter("containerName"),
	}
	if t1.containerName == "" {
		cn, _ := GetFirstContainerName(t1.ns, t1.podName)
		t1.containerName = cn
	}
	t, ok := t1.getTerminal()
	if !ok {
		resp.WriteErrorString(500, t1.key()+"没有 Exec 实例")
		return
	}
	size := &struct {
		Width  uint16
		Height uint16
	}{}
	err := req.ReadEntity(size)
	if err != nil {
		resp.WriteErrorString(500, err.Error())
		return
	}
	t.size <- &remotecommand.TerminalSize{
		Width:  size.Width,
		Height: size.Height,
	}

}

// 	ws.Route(ws.GET("/ns/{ns}/podName/{podName}/exec")
func Exec(request *restful.Request, response *restful.Response) {
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

	cancelCtx, cancel := context.WithCancel(request.Request.Context())
	eg, ctx := errgroup.WithContext(cancelCtx)

	t := &terminal{
		conn:          c,
		ns:            ns,
		podName:       podName,
		containerName: containerName,
		size:          make(chan *remotecommand.TerminalSize, 1),
	}
	t.saveTerminal()
	defer t.removeTerminal()

	t.executor(ctx, eg)

	err = eg.Wait()
	if err != nil {
		cancel()
	}
}

func (t *terminal) executor(ctx context.Context, eg *errgroup.Group) {

	eg.Go(func() error {
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
					Command:   constants.DefaultShellExecCommand,
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
	})

}
