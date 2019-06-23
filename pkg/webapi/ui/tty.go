package ui

import (
	"bytes"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/weibaohui/mesh/pkg/server"
	"io"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"net/http"
)

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

	logEvent := make(chan []byte)
	exec := executor(ns, podName, containerName)

	go func() {
		for {
			if _, reader, err := c.NextReader(); err == nil {
				var execOut bytes.Buffer
				go func() {
					err = exec.Stream(remotecommand.StreamOptions{
						Stdin:             reader,
						Stdout:            &execOut,
						Stderr:            nil,
						Tty:               true,
						TerminalSizeQueue: nil,
					})
				}()

				fmt.Println("读取一行")
				line, err := ioutil.ReadAll(&execOut)
				fmt.Println("读取", string(line))
				if err != nil {
					fmt.Println("读取", err.Error())
				}
				logEvent <- line

			}
		}
	}()

	for {
		select {
		case item, ok := <-logEvent:
			fmt.Println("logEvent接受一行")
			if !ok {

				break
			}
			if err := writeData(c, item); err != nil {
				fmt.Println("writeData", err.Error())
			}

		}
	}

	// for {
	// 	_, command, err := c.ReadMessage()
	// 	if err != nil {
	// 		log.Println("read:", err)
	// 		break
	// 	}
	// 	s := strings.TrimSpace(string(command))
	// 	fmt.Println("收到", s, strings.HasSuffix(s, "]"))
	// 	commands := strings.SplitN(s, " ", -1)
	// 	fmt.Println("split result", commands, len(commands))
	// 	for k, v := range commands {
	// 		fmt.Println("split ff ", k, v)
	// 	}
	// 	result, err := execIntoPod(ns, podName, containerName, commands)
	// 	if err != nil {
	// 		fmt.Println("execIntoPod", err.Error())
	// 	}
	// 	err = writeData(c, []byte(result))
	// 	if err != nil {
	// 		fmt.Println("write:", err)
	// 		break
	// 	}
	// }

}

func executor(ns, podName, containerName string) remotecommand.Executor {

	mctx := server.GlobalContext()
	config := server.GlobalKubeConfig()
	req := mctx.K8s.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(ns).
		SubResource("exec").
		VersionedParams(
			&v1.PodExecOptions{
				TypeMeta:  v12.TypeMeta{},
				Stdin:     true,
				Stdout:    true,
				TTY:       true,
				Container: containerName,
				Command:   []string{"/bin/sh"},
			},
			scheme.ParameterCodec,
		)
	restConfig, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		fmt.Errorf("clientcmd.BuildConfigFromFlags =%s ", err.Error())
		return nil
	}
	exec, err := remotecommand.NewSPDYExecutor(restConfig, http.MethodPost, req.URL())
	if err != nil {
		fmt.Errorf("remotecommand.NewSPDYExecutor =%s ", err.Error())
		return nil
	}
	return exec

}

func execIntoPodReader(in io.Reader, logEvent chan []byte, ns, podName, containerName string) {

	var execOut bytes.Buffer
	mctx := server.GlobalContext()
	config := server.GlobalKubeConfig()
	req := mctx.K8s.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(ns).
		SubResource("exec").
		VersionedParams(
			&v1.PodExecOptions{
				TypeMeta:  v12.TypeMeta{},
				Stdin:     true,
				Stdout:    true,
				TTY:       true,
				Container: containerName,
				Command:   []string{"/bin/sh"},
			},
			scheme.ParameterCodec,
		)
	restConfig, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		fmt.Errorf("clientcmd.BuildConfigFromFlags =%s ", err.Error())
		return
	}
	exec, err := remotecommand.NewSPDYExecutor(restConfig, http.MethodPost, req.URL())
	if err != nil {
		fmt.Errorf("remotecommand.NewSPDYExecutor =%s ", err.Error())
		return
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  in,
		Stdout: &execOut,
		Tty:    true,
	})

	if err != nil {
		fmt.Errorf("could not execute: %v", err)
		return
	}

	for {
		line, err := execOut.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		logEvent <- line

	}

}

func execIntoPod(ns, podName, containerName string, commands []string) (string, error) {

	var execOut bytes.Buffer
	var execErr bytes.Buffer
	mctx := server.GlobalContext()
	config := server.GlobalKubeConfig()
	req := mctx.K8s.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(ns).
		SubResource("exec").
		VersionedParams(
			&v1.PodExecOptions{
				TypeMeta:  v12.TypeMeta{},
				Stdin:     false,
				Stdout:    true,
				Stderr:    true,
				TTY:       true,
				Container: containerName,
				Command:   commands,
			},
			scheme.ParameterCodec,
		)
	restConfig, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		fmt.Errorf("clientcmd.BuildConfigFromFlags =%s ", err.Error())
		return "", err
	}
	exec, err := remotecommand.NewSPDYExecutor(restConfig, http.MethodPost, req.URL())
	if err != nil {
		fmt.Errorf("remotecommand.NewSPDYExecutor =%s ", err.Error())
		return "", err
	}

	err = exec.Stream(remotecommand.StreamOptions{
		// Stdin:             strings.NewReader(commands),
		Stdout: &execOut,
		Stderr: &execErr,
		Tty:    true,
	})

	if err != nil {
		fmt.Errorf("could not execute: %v", err)
		return "", err
	}

	if execErr.Len() > 0 {
		fmt.Errorf("stderr: %v", execErr.String())
		return "", err
	}
	for {
		line, err := execOut.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		fmt.Println(string(line))

	}

	return execOut.String(), nil

}
