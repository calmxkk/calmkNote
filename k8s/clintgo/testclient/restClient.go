package testclient

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func testRESTClient() {
	cfg, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
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
	restClient.Get().Namespace("default").Resource("pod").Do(context.TODO()).Into(result)
	for _, r := range result.Items {
		fmt.Println(r.Name)
	}
}
