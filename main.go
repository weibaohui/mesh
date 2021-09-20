//go:generate go run pkg/codegen/cleanup/main.go
//go:generate go run vendor/github.com/jteeuwen/go-bindata/go-bindata/AppendSliceValue.go vendor/github.com/jteeuwen/go-bindata/go-bindata/main.go vendor/github.com/jteeuwen/go-bindata/go-bindata/version.go -o ./stacks/bindata.go -ignore bindata.go -pkg stacks -modtime 1557785965 -mode 0644 ./stacks/
//go:generate go fmt stacks/bindata.go
//go:generate go run pkg/codegen/main.go

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/rancher/wrangler/pkg/signals"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/weibaohui/mesh/pkg/constants"
	"github.com/weibaohui/mesh/pkg/server"
	"github.com/weibaohui/mesh/pkg/version"
	"github.com/weibaohui/mesh/pkg/webapi"
	"k8s.io/klog"
)

var (
	debug      bool
	kubeconfig string
	namespace  string
)

func main() {
	app := cli.NewApp()
	app.Name = "mesh-controller"
	app.Version = fmt.Sprintf("%s (%s)", version.Version, version.GitCommit)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "kubeconfig",
			EnvVar:      "KUBECONFIG",
			Destination: &kubeconfig,
		},
		cli.StringFlag{
			Name:        "namespace",
			EnvVar:      "MESH_NAMESPACE",
			Value:       "mesh-system",
			Destination: &namespace,
		},
		cli.BoolFlag{
			Name:        "debug",
			EnvVar:      "MESH_DEBUG",
			Destination: &debug,
		},
		cli.StringFlag{
			Name:        "service-cidr",
			Usage:       "Manually specify cluster IP CIDR for envoy",
			EnvVar:      "SERVICE_CIDR",
			Value:       "10.43.0.0/16",
			Destination: &constants.ServiceCidr,
		},
	}
	app.Action = run

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {

	debug = true
	if debug {
		setupDebugLogging()
		logrus.SetLevel(logrus.DebugLevel)
	}

	ctx := signals.SetupSignalHandler(context.Background())

	go webapi.Start(ctx)
	if err := server.Startup(ctx, namespace, kubeconfig); err != nil {
		return err
	}
	builder := webapi.WorkLoadBuilder{}
	builder.Load("default", "weave-scope-app").String()
	return nil
}

func setupDebugLogging() {
	flag.Set("alsologtostderr", "true")
	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(klogFlags)
}
