package api

import v1 "k8s.io/api/core/v1"

type Containers []Container

type Container struct {
	Name   string
	Pod    Pod
	IsInit bool
}

func NewContainer(pod v1.Pod, container v1.Container, isInit bool) Container {
	return Container{
		Name:   container.Name,
		Pod:    NewPod(pod),
		IsInit: isInit,
	}
}

type Pod struct {
	Name      string
	Namespace string
	Phase     string
}

func NewPod(pod v1.Pod) Pod {
	return Pod{
		Name:      pod.Name,
		Namespace: pod.Namespace,
		Phase:     string(pod.Status.Phase),
	}
}
