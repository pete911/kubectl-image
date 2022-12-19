package out

import (
	"fmt"
	"github.com/pete911/kubectl-image/pkg/api"
)

func PrintRegistries(registries api.Registries, nodes api.Nodes) {

	for _, registry := range registries {
		printRegistry(registry, nodes)
	}
}

func printRegistry(registry api.Registry, nodes api.Nodes) {

	registryName := registry.Name
	if registryName == "" {
		registryName = "-"
	}
	fmt.Printf("registry: %s\n", registryName)
	for _, repository := range registry.ListRepositories() {
		printRepository(repository, nodes)
	}
}

func printRepository(repository api.Repository, nodes api.Nodes) {

	fmt.Printf("  %s\n", repository.Name)
	for _, tagID := range repository.ListTagIDs() {
		printTagID(tagID, nodes)
	}
}

func printTagID(tagDigest api.TagID, nodes api.Nodes) {

	sizeStr := ""
	if size := nodes.GetSizeBytes(tagDigest.ImageName); size != 0 {
		sizeMb := float64(size) / 1000000
		sizeStr = fmt.Sprintf("\tSize: %.2fMB", sizeMb)
	}
	fmt.Printf("    Tag/ID: %s%s\n", tagDigest.Name, sizeStr)
	for _, id := range tagDigest.ListIDs() {
		printID(id)
	}
}

func printID(id api.ID) {

	fmt.Printf("    ID:     %s\n", id.Name)
	for _, container := range id.ListContainers() {
		printContainer(container)
	}
}

func printContainer(container api.Container) {

	containerKey := "[container]"
	if container.IsInit {
		containerKey = "[init-container]"
	}
	fmt.Printf("            [namespace] %s %s %s [pod] %s [pod-phase] %s\n",
		container.Pod.Namespace, containerKey, container.Name, container.Pod.Name, container.Pod.Phase)
}
