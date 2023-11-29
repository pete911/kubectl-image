package api

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"strings"
	"time"
)

type Containers []Container

type Container struct {
	Name         string
	Pod          Pod
	IsInit       bool
	ImageName    ImageName
	State        string // waiting, running or terminated
	Message      string // reason and message if container is waiting, start time if running, ...
	RestartCount int
}

func NewContainer(pod v1.Pod, container v1.Container, isInit bool) Container {
	containerStatus := getContainerStatus(pod, container.Name, isInit)
	state, message := getContainerStateAndMessage(containerStatus)
	// set both, image with tag and id as well (set by kube)
	imageName := ParseImageName(container.Image)
	imageName.ID = ParseImageID(containerStatus.ImageID)

	return Container{
		Name:         container.Name,
		Pod:          NewPod(pod),
		IsInit:       isInit,
		ImageName:    imageName,
		State:        state,
		Message:      message,
		RestartCount: int(containerStatus.RestartCount),
	}
}

func getContainerStateAndMessage(containerStatus v1.ContainerStatus) (string, string) {
	state := containerStatus.State
	if state.Waiting != nil {
		return "waiting", toMessage(state.Waiting.Reason, state.Waiting.Message)
	}
	if state.Running != nil {
		return "running", fmt.Sprintf("started at %s", state.Running.StartedAt.Format(time.RFC3339))
	}
	if state.Terminated != nil {
		return "terminated", toMessage(fmt.Sprintf("exit code: %d", state.Terminated.ExitCode), state.Terminated.Reason, state.Terminated.Message)
	}
	return "", ""
}

func toMessage(in ...string) string {
	return strings.TrimSpace(strings.Join(in, " "))
}

func getContainerStatus(pod v1.Pod, containerName string, isInit bool) v1.ContainerStatus {
	for _, status := range getContainerStatuses(pod, isInit) {
		if status.Name == containerName {
			return status
		}
	}
	return v1.ContainerStatus{}
}

func getContainerStatuses(pod v1.Pod, isInit bool) []v1.ContainerStatus {
	if isInit {
		return pod.Status.InitContainerStatuses
	}
	return pod.Status.ContainerStatuses
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
