package main

import (
	"fmt"
	"io"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubernetes/pkg/kubectl/util/term"
	"net/http"
	"os"
)

func mains()  {
	http.HandleFunc("/tty", func(writer http.ResponseWriter, request *http.Request) {
		getclient(request.Body,writer)
	})
	http.ListenAndServe(":6666",nil)
}
func getclient(in io.Reader, out io.Writer) {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/baohui/.kube/config")
	if err != nil {
		fmt.Println(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
	}
	req := clientset.CoreV1().RESTClient().Get().
		Resource("pods").
		Name("busybox-57f498dd4f-mh58v").
		Namespace("default").
		SubResource("log")

	req.VersionedParams(
		&v1.PodExecOptions{
			Container: "busybox",
			Command:   []string{"ls"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		},
		scheme.ParameterCodec,

	)

	closer, err := req.Stream()
	bytes, err := ioutil.ReadAll(closer)
	fmt.Println(string(bytes))

	out.Write(bytes)
	fmt.Println("1")
	executor, err := remotecommand.NewSPDYExecutor(
		config, http.MethodPost, req.URL(),
	)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("2")
	tty := term.TTY{}
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             os.Stdin,
		Stdout:            out,
		Stderr:            os.Stderr,
		Tty:               true,
		TerminalSizeQueue: tty.MonitorSize(),
	})
	fmt.Println("3")

	if err != nil {
		fmt.Println(err.Error())
	}

}
