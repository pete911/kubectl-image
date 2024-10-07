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
	Namespace string
	coreV1    corev1.CoreV1Interface
}

func NewClient(kubeconfigPath, contextName, namespace string) (Client, error) {
	clientConfig := newClientConfig(kubeconfigPath, contextName)
	if namespace == "" {
		if ns, _, err := clientConfig.Namespace(); err == nil {
			namespace = ns
		}
		if namespace == "" {
			namespace = "default"
		}
	}

	config, err := clientConfig.ClientConfig()
	if err != nil {
		return Client{}, err
	}
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		return Client{}, err
	}
	return Client{
		Namespace: namespace,
		coreV1:    cs.CoreV1(),
	}, nil
}

func newClientConfig(kubeconfigPath, contextName string) clientcmd.ClientConfig {
	if contextName == "" {
		return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
			&clientcmd.ConfigOverrides{},
		)
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: contextName,
		},
	)
}

func (c Client) ListRegistries(allNamespaces bool) (Registries, error) {
	nodes, err := c.listNodes()
	if err != nil {
		return nil, fmt.Errorf("list nodes: %w", err)
	}

	if allNamespaces {
		pods, err := c.getAllPods()
		if err != nil {
			return nil, fmt.Errorf("get images: %w", err)
		}
		return NewRegistries(nodes, pods), nil
	}

	pods, err := c.getPods(c.Namespace)
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
