package api

import (
	"context"
	"fmt"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"time"
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

func (c Client) ListRegistries(namespace string) (Registries, error) {
	nodes, err := c.listNodes()
	if err != nil {
		return nil, fmt.Errorf("list nodes: %w", err)
	}

	if namespace == "" {
		pods, err := c.getAllPods()
		if err != nil {
			return nil, fmt.Errorf("get images: %w", err)
		}
		return NewRegistries(nodes, pods), nil
	}

	pods, err := c.getPods(namespace)
	if err != nil {
		return nil, fmt.Errorf("get images: %w", err)
	}
	return NewRegistries(nodes, pods), nil
}

func (c Client) getAllPods() ([]v1.Pod, error) {

	namespaces, err := c.getNamespaces()
	if err != nil {
		return nil, fmt.Errorf("get namespaces: %w", err)
	}

	var pods []v1.Pod
	for _, namespace := range namespaces {
		p, err := c.getPods(namespace.Name)
		if err != nil {
			return nil, err
		}
		pods = append(pods, p...)
	}
	return pods, nil
}

func (c Client) getPods(namespace string) ([]v1.Pod, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	podList, err := c.coreV1.Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return podList.Items, nil
}

func (c Client) getNamespaces() ([]v1.Namespace, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	namespaceList, err := c.coreV1.Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("get namespaces: %w", err)
	}
	return namespaceList.Items, nil
}

func (c Client) listNodes() (map[string]Node, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	nodeList, err := c.coreV1.Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("get nodes: %w", err)
	}
	return NewNodes(nodeList.Items), nil
}
