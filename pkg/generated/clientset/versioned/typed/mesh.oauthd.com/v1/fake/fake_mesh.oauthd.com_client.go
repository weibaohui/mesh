/*

 */

// Code generated by ___go_build_main_go. DO NOT EDIT.

package fake

import (
	v1 "github.com/weibaohui/mesh/pkg/generated/clientset/versioned/typed/mesh.oauthd.com/v1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeMeshV1 struct {
	*testing.Fake
}

func (c *FakeMeshV1) Apps(namespace string) v1.AppInterface {
	return &FakeApps{c, namespace}
}

func (c *FakeMeshV1) ClusterDomains(namespace string) v1.ClusterDomainInterface {
	return &FakeClusterDomains{c, namespace}
}

func (c *FakeMeshV1) ExternalServices(namespace string) v1.ExternalServiceInterface {
	return &FakeExternalServices{c, namespace}
}

func (c *FakeMeshV1) Routers(namespace string) v1.RouterInterface {
	return &FakeRouters{c, namespace}
}

func (c *FakeMeshV1) Services(namespace string) v1.ServiceInterface {
	return &FakeServices{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeMeshV1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
