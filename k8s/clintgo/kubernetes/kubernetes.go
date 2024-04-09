package kubernetes

import (
	"clinetgo/model"
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
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

type K8sClient struct {
	config *rest.Config
	client *kubernetes.Clientset
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

func (k *sKubernetesCluster) Client(ctx context.Context, in *model.Cluster) (*K8sClient, error) {
	cfg, err := k.Config(ctx, in)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	client, _ := kubernetes.NewForConfig(cfg)

	return &K8sClient{
		config: cfg,
		client: client,
	}, nil
}

func (k *K8sClient) Version() (*version.Info, error) {
	return k.client.ServerVersion()
}

func (k *K8sClient) Ping(ctx context.Context) error {
	_, err := k.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	k.client.AuthorizationV1().SelfSubjectAccessReviews()
	if err != nil {
		return err
	}

	return nil
}

func (k *K8sClient) GetUserNamespaceNames(ctx context.Context) error {
	namespaces, err := k.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil
	}
	for _, ns := range namespaces.Items {
		fmt.Println(ns.Name)
	}
	return nil
}

func (k *K8sClient) GetPodByNamespace(ctx context.Context, namespace string) error {
	pods, err := k.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, pod := range pods.Items {
		fmt.Printf("%s\n", pod.Name)
	}

	deployments, err := k.client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	fmt.Println()

	for _, pod := range deployments.Items {
		fmt.Printf("%s\n", pod.Name)
	}

	return nil
}

func (k *K8sClient) CreatePod(ctx context.Context, in *model.Pod) (*corev1.Pod, error) {
	if in == nil {
		return nil, fmt.Errorf("%s", "no input data for create pod")
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      in.Metadata.Name,
			Labels:    in.Metadata.Lables,
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            in.PodSpec.Containers[0].Name,
					Image:           in.PodSpec.Containers[0].Image,
					ImagePullPolicy: corev1.PullPolicy(in.PodSpec.Containers[0].ImagePullPolicy),
				},
			},
		},
	}

	get_pod, err := k.client.CoreV1().Pods("default").Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return get_pod, nil
}
