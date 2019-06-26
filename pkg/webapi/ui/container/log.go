package container

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/gorilla/websocket"
	"github.com/weibaohui/mesh/pkg/server"
	"golang.org/x/sync/errgroup"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 	ws.Route(ws.GET("/ns/{ns}/podName/{podName}/log")
func Log(request *restful.Request, response *restful.Response) {
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
	readerGroup, ctx := errgroup.WithContext(cancelCtx)

	go func() {
		for {
			if _, _, err := c.NextReader(); err != nil {
				cancel()
				c.Close()
				break
			}
		}
	}()
	logEvent := make(chan []byte)

	ReadLog(ctx, readerGroup, logEvent, ns, podName, containerName)

	go func() {
		readerGroup.Wait()
		close(logEvent)
	}()
	done := false
	for !done {
		select {
		case item, ok := <-logEvent:
			if !ok {
				done = true
				break
			}
			if err := writeData(c, item); err != nil {
				cancel()
			}

		}
	}

}

func writeData(c *websocket.Conn, buf []byte) error {
	messageWriter, err := c.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	if _, err := messageWriter.Write(buf); err != nil {
		return err
	}
	return messageWriter.Close()
}
func GetFirstContainerName(ns string, podName string) (string, error) {
	mctx := server.GlobalContext()
	pod, err := mctx.Core.Core().V1().Pod().Cache().Get(ns, podName)
	if err != nil {
		return "", err
	}
	if len(pod.Spec.Containers) == 0 {
		return "", errors.New("没有容器")
	}
	return pod.Spec.Containers[0].Name, nil
}

func ReadLog(ctx context.Context, eg *errgroup.Group, logEvent chan []byte, ns, podName, containerName string) {
	eg.Go(func() error {
		mctx := server.GlobalContext()
		req := mctx.K8s.CoreV1().RESTClient().Get().
			Resource("pods").
			Name(podName).
			Namespace(ns).
			SubResource("log").
			VersionedParams(
				&v1.PodLogOptions{
					Container: containerName,
					Follow:    true,
				},
				scheme.ParameterCodec,
			)
		readCloser, err := req.Stream()
		if err != nil {
			fmt.Println("Stream", err)
			return err
		}
		reader := bufio.NewReader(readCloser)
		for {
			line, _, err := reader.ReadLine()
			logEvent <- line
			if err == io.EOF {
				break
			}
		}
		return nil
	})

}
