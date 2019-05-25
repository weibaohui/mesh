package serviceset

import (
	v1 "github.com/weibaohui/mesh/pkg/apis/mesh.oauthd.com/v1"
)

type Services map[string]*ServiceSet

func (s Services) List() []*v1.Service {
	var result []*v1.Service
	for _, v := range s {
		result = append(result, v.Revisions...)
	}
	return result
}

type ServiceSet struct {
	Revisions []*v1.Service
}
