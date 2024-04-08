package kubernetes

import (
	"clinetgo/model/cluster"

	"k8s.io/client-go/rest"
)

type K8s struct {
	*cluster.Cluster
}

func (k *K8s) Config() (*rest.Config, error) {
	if k.Spec.Local {
		return rest.InClusterConfig()
	}

	if k.Spec.Connect.Direction == "forward" {
		kubeConf := &rest.Config{
			Host: k.Spec.Connect.Forward.ApiServer,
		}
	}
}
