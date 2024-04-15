package testclient

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

func TestDiscoveryClient() {
	cfg, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		fmt.Println(err)
		return
	}

	discoveryclient, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, apiResoure, err := discoveryclient.ServerGroupsAndResources()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, list := range apiResoure {
		gv, err := schema.ParseGroupVersion(list.GroupVersion)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, res := range list.APIResources {
			fmt.Println(res.Name, gv.Group)
		}
	}
}
