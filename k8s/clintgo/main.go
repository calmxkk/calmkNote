package main

import (
	"clinetgo/kubernetes"
	"clinetgo/model"
	"context"
	"fmt"
	"os"

	"github.com/gogf/gf/v2/os/gfile"
)

func main() {
	TestCluster()

}

func TestCluster() {
	configfilebytes := gfile.GetBytes(os.Getenv("HOME") + "/.kube/config")
	cluster := model.Cluster{
		Spec: model.Spec{
			Connect: model.Connect{Direction: "forward"},
			Authentication: model.Authentication{
				Mode:              "configfile",
				ConfigFileContent: configfilebytes},
		},
		Namespace: "kube-system",
	}
	ctx := context.Background()
	k8s := kubernetes.NewKubernetesCluster()
	version, err := k8s.Version(ctx, &cluster)
	if err != nil {
		return
	}
	fmt.Println(version)

	err = k8s.GetPodByNamespace(ctx, &cluster)
	fmt.Println(err)
}
