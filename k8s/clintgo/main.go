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
	k8s, _ := kubernetes.NewKubernetesCluster().Client(ctx, &cluster)
	version, err := k8s.Version()
	if err != nil {
		return
	}
	fmt.Println(version)

	err = k8s.GetPodByNamespace(ctx, cluster.Namespace)
	fmt.Println(err)
	err = k8s.GetUserNamespaceNames(ctx)
	fmt.Println(err)

}
