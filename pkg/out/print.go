package out

import (
	"fmt"
	"github.com/pete911/kubectl-image/pkg/api"
)

func Print(registries api.Registries) {

	for _, registry := range registries {
		printRegistry(registry, "")
	}
}

func printRegistry(registry api.Registry, prefix string) {

	registryName := registry.Name
	if registryName == "" {
		registryName = "-"
	}
	fmt.Printf("%sregistry: %s\n", prefix, registryName)
	for _, repository := range registry.ListRepositories() {
		printRepository(repository, "  ")
	}
}

func printRepository(repository api.Repository, prefix string) {

	fmt.Printf("%s%s\n", prefix, repository.Name)
	for _, tagID := range repository.ListTagIDs() {
		printTagID(tagID, "    ")
	}
}

func printTagID(tagDigest api.TagID, prefix string) {

	fmt.Printf("%sTag/ID: %s\n", prefix, tagDigest.Name)
	for _, id := range tagDigest.ListIDs() {
		printID(id, "    ")
	}
}

func printID(id api.ID, prefix string) {

	fmt.Printf("%sID:     %s\n", prefix, id.Name)
	for _, container := range id.ListContainers() {
		printContainer(container, "            ")
	}
}

func printContainer(container api.Container, prefix string) {

	containerKey := "[container]"
	if container.IsInit {
		containerKey = "[init-container]"
	}
	fmt.Printf("%s[namespace] %s %s %s [pod] %s [pod-phase] %s\n",
		prefix, container.Pod.Namespace, containerKey, container.Name, container.Pod.Name, container.Pod.Phase)
}
