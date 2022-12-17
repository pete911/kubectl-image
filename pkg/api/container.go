package api

import v1 "k8s.io/api/core/v1"

type Containers []Container

type Container struct {
	Name    string
	Pod     Pod
	IsInit  bool
	ImageID string
}

func NewContainer(pod v1.Pod, container v1.Container, isInit bool) Container {
	var imageID string
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if container.Name == containerStatus.Name {
			// sometimes the image ID contains only ID, sometimes the whole image
			imageID = ParseImageID(containerStatus.ImageID)
		}
	}

	return Container{
		Name:    container.Name,
		Pod:     NewPod(pod),
		IsInit:  isInit,
		ImageID: imageID,
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
