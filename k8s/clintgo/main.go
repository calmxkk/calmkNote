package main

import (
	"clinetgo/kubernetes"
	"clinetgo/model"
	"context"
	"fmt"

	"github.com/gogf/gf/v2/os/gfile"
)

func main() {
	TestCluster()

}

func TestCluster() {
	configfilebytes := gfile.GetBytes("/root/.kube/config")
	cluster := model.Cluster{
		Spec: model.Spec{
			Connect: model.Connect{Direction: "forward"},
			Authentication: model.Authentication{
				Mode:              "configfile",
				ConfigFileContent: configfilebytes},
		},
	}
	ctx := context.Background()
	k8s := kubernetes.NewKubernetesCluster()
	version, err := k8s.Version(ctx, &cluster)
	if err != nil {
		return
	}
	fmt.Println(version)

	err = k8s.Ping(ctx, &cluster)
	fmt.Println(err)
}
