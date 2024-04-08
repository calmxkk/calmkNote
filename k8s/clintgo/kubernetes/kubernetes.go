package kubernetes

import (
	"clinetgo/model"
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type sKubernetesCluster struct{}

func NewKubernetesCluster() *sKubernetesCluster {
	return &sKubernetesCluster{}
}

func (k *sKubernetesCluster) Config(ctx context.Context, in *model.Cluster) (*rest.Config, error) {
	if in.Spec.Local {
		return rest.InClusterConfig()
	}

	if in.Spec.Connect.Direction == "forward" {
		kubeConf := &rest.Config{
			Host: in.Spec.Connect.Forward.ApiServer,
		}

		if len(in.CaCertificate.CertData) > 0 {
			kubeConf.CAData = in.CaCertificate.CertData
		} else {
			kubeConf.Insecure = true
		}

		switch strings.ToLower(in.Spec.Authentication.Mode) {
		case "bearer":
			kubeConf.BearerToken = in.Spec.Authentication.BearerToken
		case "certificate":
			kubeConf.TLSClientConfig.CertData = in.Spec.Authentication.Certificate.CertData
			kubeConf.TLSClientConfig.KeyData = in.Spec.Authentication.Certificate.KeyData
		case "configfile":
			cfg, err := clientcmd.BuildConfigFromKubeconfigGetter("", func() (*clientcmdapi.Config, error) {
				return clientcmd.Load(in.Spec.Authentication.ConfigFileContent)
			})
			if err != nil {
				return nil, err
			}
			kubeConf = cfg
		}
		return kubeConf, nil
	}

	return nil, nil
}

func (k *sKubernetesCluster) Client(ctx context.Context, in *model.Cluster) (*kubernetes.Clientset, error) {
	cfg, err := k.Config(ctx, in)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return kubernetes.NewForConfig(cfg)
}

func (k *sKubernetesCluster) Version(ctx context.Context, in *model.Cluster) (*version.Info, error) {
	client, err := k.Client(ctx, in)
	if err != nil {
		return nil, err
	}
	return client.ServerVersion()
}

func (k *sKubernetesCluster) Ping(ctx context.Context, in *model.Cluster) error {
	client, err := k.Client(ctx, in)
	if err != nil {
		return err
	}
	_, err = client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	client.AuthorizationV1().SelfSubjectAccessReviews()
	if err != nil {
		return err
	}

	return nil
}

// func (k *sKubernetesCluster) GetUserNamespaceNames(ctx context.Context, in *model.Cluster)

func (k *sKubernetesCluster) GetPodByNamespace(ctx context.Context, in *model.Cluster) error {
	client, err := k.Client(ctx, in)
	if err != nil {
		return err
	}

	pods, err := client.CoreV1().Pods(in.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, pod := range pods.Items {
		fmt.Printf("%s\n", pod.Name)
	}

	deployments, err := client.AppsV1().Deployments(in.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	fmt.Println()

	for _, pod := range deployments.Items {
		fmt.Printf("%s\n", pod.Name)
	}

	return nil
}
