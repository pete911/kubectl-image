package api

import (
	"k8s.io/api/core/v1"
	"sort"
)

// Registries map where key is registry name (e.g. gcr.io) and value is Registry struct
type Registries map[string]Registry

func NewRegistries(pods []v1.Pod) Registries {
	registries := Registries{}
	for _, p := range pods {
		for _, c := range p.Spec.Containers {
			addToRegistries(registries, p, c, false)
		}
		for _, c := range p.Spec.InitContainers {
			addToRegistries(registries, p, c, true)
		}
	}
	return registries
}

func addToRegistries(registries Registries, p v1.Pod, c v1.Container, isInit bool) {
	container := NewContainer(p, c, isInit)
	if _, ok := registries[container.ImageName.Registry]; !ok {
		registries[container.ImageName.Registry] = newRegistry(container)
	}
	registries[container.ImageName.Registry].addRepository(container)
}

func (r Registries) List() []Registry {
	var out []Registry
	for _, v := range r {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[j].Name > out[i].Name
	})
	return out
}

// --- registry ---

// Registry (e.g. gcr.io) 'bucket' of repositories (and images)
type Registry struct {
	Name         string
	repositories map[string]Repository
}

func newRegistry(container Container) Registry {
	return Registry{
		Name:         container.ImageName.Registry,
		repositories: map[string]Repository{},
	}
}

func (r Registry) ListRepositories() []Repository {
	var out []Repository
	for _, v := range r.repositories {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[j].Name > out[i].Name
	})
	return out
}

func (r Registry) addRepository(container Container) {
	if _, ok := r.repositories[container.ImageName.Repository]; !ok {
		r.repositories[container.ImageName.Repository] = newRepository(container)
	}
	r.repositories[container.ImageName.Repository].addID(container)
}

// --- repository ---

// Repository (image name without registry and tag/id e.g. jacksontj/promxy) 'bucket' of tags/ids
type Repository struct {
	Name string
	IDs  map[string]ID
}

func newRepository(container Container) Repository {
	return Repository{
		Name: container.ImageName.Repository,
		IDs:  map[string]ID{},
	}
}

func (r Repository) ListIDs() []ID {
	var out []ID
	for _, v := range r.IDs {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[j].Name > out[i].Name
	})
	return out
}

func (r Repository) addID(container Container) {
	if _, ok := r.IDs[container.ImageName.ID]; !ok {
		r.IDs[container.ImageName.ID] = newID(container.ImageName.ID)
	}
	r.IDs[container.ImageName.ID].addContainer(container)
}

// --- container digest/id ---

// ID is container image digest/ID, it comes from container status after image is pulled
type ID struct {
	Name       string
	tags       map[string]struct{}
	containers map[containerKey]Container
}

type containerKey struct {
	podName       string
	podNamespace  string
	containerName string
}

func newID(name string) ID {
	return ID{Name: name, tags: map[string]struct{}{}, containers: map[containerKey]Container{}}
}

func (i ID) addContainer(container Container) {
	key := containerKey{
		podName:       container.Pod.Name,
		podNamespace:  container.Pod.Namespace,
		containerName: container.Name,
	}
	if container.ImageName.Tag != "" {
		i.tags[container.ImageName.Tag] = struct{}{}
	}
	i.containers[key] = container
}

func (i ID) ListContainers() []Container {
	var out []Container
	for _, v := range i.containers {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[j].Name > out[i].Name
	})
	return out
}

func (i ID) ListTags() []string {
	var out []string
	for k := range i.tags {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
