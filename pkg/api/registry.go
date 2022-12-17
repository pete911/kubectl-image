package api

import (
	"k8s.io/api/core/v1"
	"sort"
)

type Registries map[string]Registry

func NewRegistries(pods []v1.Pod) Registries {
	registries := Registries{}
	for _, p := range pods {
		for _, c := range p.Spec.Containers {
			addToRegistries(registries, p, c, false)
		}
		for _, c := range p.Spec.InitContainers {
			addToRegistries(registries, p, c, false)
		}
	}
	return registries
}

func addToRegistries(registries Registries, p v1.Pod, c v1.Container, isInit bool) {
	container := NewContainer(p, c, isInit)
	imageName := ParseImageName(c.Image)
	if _, ok := registries[imageName.Registry]; !ok {
		registries[imageName.Registry] = newRegistry(imageName)
	}
	registries[imageName.Registry].addRepository(imageName, container)
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

type Registry struct {
	Name         string
	repositories map[string]Repository
}

func newRegistry(imageName ImageName) Registry {
	return Registry{Name: imageName.Registry, repositories: map[string]Repository{}}
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

func (r Registry) addRepository(imageName ImageName, container Container) {
	if _, ok := r.repositories[imageName.Repository]; !ok {
		r.repositories[imageName.Repository] = newRepository(imageName)
	}
	r.repositories[imageName.Repository].addTagDigest(imageName, container)
}

// --- repository ---

type Repository struct {
	Name       string
	tagDigests map[string]TagDigest
}

func newRepository(imageName ImageName) Repository {
	return Repository{Name: imageName.Repository, tagDigests: map[string]TagDigest{}}
}

func (r Repository) ListTagDigests() []TagDigest {
	var out []TagDigest
	for _, v := range r.tagDigests {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[j].Name > out[i].Name
	})
	return out
}

func (r Repository) addTagDigest(imageName ImageName, container Container) {
	if _, ok := r.tagDigests[imageName.TagDigest()]; !ok {
		r.tagDigests[imageName.TagDigest()] = newTagDigest(imageName.TagDigest())
	}
	r.tagDigests[imageName.TagDigest()].addContainer(container)
}

// --- tag or digest ---

type containerKey struct {
	podName       string
	podNamespace  string
	containerName string
}

type TagDigest struct {
	Name       string
	containers map[containerKey]Container
}

func newTagDigest(name string) TagDigest {
	return TagDigest{Name: name, containers: map[containerKey]Container{}}
}

func (t TagDigest) addContainer(container Container) {
	key := containerKey{
		podName:       container.Pod.Name,
		podNamespace:  container.Pod.Namespace,
		containerName: container.Name,
	}
	t.containers[key] = container
}

func (t TagDigest) ListContainers() []Container {
	var out []Container
	for _, v := range t.containers {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[j].Name > out[i].Name
	})
	return out
}
