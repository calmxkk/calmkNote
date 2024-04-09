package testclient

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func testDynamicClient() {
	cfg, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		fmt.Println(err)
		return
	}

	dyclient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	gvr := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}

	unstructDaya, err := dyclient.Resource(gvr).Namespace("default").List(context.TODO(), metav1.ListOptions{})

	podList := &corev1.PodList{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(
		unstructDaya.UnstructuredContent(),
		podList)

	for _, p := range podList.Items {
		fmt.Println(p.Name, p.Labels)
	}
}
