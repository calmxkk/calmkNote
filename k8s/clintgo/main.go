package main

import (
	"clinetgo/kubernetes"
	"clinetgo/model"
	"context"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"

	"github.com/gogf/gf/v2/os/gfile"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	testRESTclient()

}

func testRESTclient() {
	cfg, err := clientcmd.BuildConfigFromFlags("". "/root/.kube/config")
	if err != nil {
		fmt.Println(err)
		return
	}

	restClient, err := rest.RESTClientFor(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	result := &corev1.PodList{}
	restClient.Get().Namespace("default").Resource("pod").Do().Into(result)
	for _, r := range result.Items {
		fmt.Println(r.Name)
	}
}

func TestPod() {
	client := getClient()

	pod := &model.Pod{
		Metadata: model.Metadata{
			Name: "nginx",
			Lables: map[string]string{
				"app": "calmk",
			},
		},
		PodSpec: model.PodSpec{
			Containers: []model.Container{
				{
					Name:            "nginx",
					Image:           "ngins:latest",
					ImagePullPolicy: "IfNotPresent",
				},
			},
		},
	}

	_, _ = client.CreatePod(context.Background(), pod)
}
func TestCluster() {
	client := getClient()



	version, err := client.Version()
	if err != nil {
		return
	}
	fmt.Println(version)

	// err = client.GetPodByNamespace(context.Background(), cluster.Namespace)
	// fmt.Println(err)
	// err = client.GetUserNamespaceNames(ctx)
	// fmt.Println(err)
}

func getClient() *kubernetes.K8sClient {
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
	return k8s
}
