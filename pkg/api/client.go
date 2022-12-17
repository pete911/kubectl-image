package api

import (
	"context"
	"fmt"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	coreV1 corev1.CoreV1Interface
}

func NewClient(kubeconfigPath string) (Client, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return Client{}, err
	}
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		return Client{}, err
	}
	return Client{coreV1: cs.CoreV1()}, nil
}

func (c Client) ListRegistries(ctx context.Context) (Registries, error) {

	pods, err := c.getAllPods(ctx, "", "")
	if err != nil {
		return nil, fmt.Errorf("get images: %w", err)
	}
	return NewRegistries(pods), nil
}

func (c Client) getAllPods(ctx context.Context, labelSelector, fieldSelector string) ([]v1.Pod, error) {

	namespaces, err := c.getNamespaces(ctx)
	if err != nil {
		return nil, fmt.Errorf("get namespaces: %w", err)
	}

	var pods []v1.Pod
	for _, namespace := range namespaces {
		p, err := c.getPods(ctx, namespace.Name, labelSelector, fieldSelector)
		if err != nil {
			return nil, err
		}
		pods = append(pods, p...)
	}
	return pods, nil
}

func (c Client) getPods(ctx context.Context, namespace, labelSelector, fieldSelector string) ([]v1.Pod, error) {

	podList, err := c.coreV1.Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: labelSelector, FieldSelector: fieldSelector})
	if err != nil {
		return nil, err
	}
	return podList.Items, nil
}

func (c Client) getNamespaces(ctx context.Context) ([]v1.Namespace, error) {

	namespaceList, err := c.coreV1.Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("get namespaces: %w", err)
	}
	return namespaceList.Items, nil
}
