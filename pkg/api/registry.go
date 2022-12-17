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
	r.repositories[imageName.Repository].addTagID(imageName, container)
}

// --- repository ---

type Repository struct {
	Name   string
	tagIDs map[string]TagID
}

func newRepository(imageName ImageName) Repository {
	return Repository{Name: imageName.Repository, tagIDs: map[string]TagID{}}
}

func (r Repository) ListTagIDs() []TagID {
	var out []TagID
	for _, v := range r.tagIDs {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[j].Name > out[i].Name
	})
	return out
}

func (r Repository) addTagID(imageName ImageName, container Container) {
	if _, ok := r.tagIDs[imageName.TagOrID()]; !ok {
		r.tagIDs[imageName.TagOrID()] = newTagID(imageName.TagOrID())
	}
	r.tagIDs[imageName.TagOrID()].addID(container)
}

// --- image tag or digest/id ---

// TagID is image tag or id, it comes from container
type TagID struct {
	Name string
	iDs  map[string]ID
}

func newTagID(name string) TagID {
	return TagID{Name: name, iDs: map[string]ID{}}
}

func (t TagID) ListIDs() []ID {
	var out []ID
	for _, v := range t.iDs {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[j].Name > out[i].Name
	})
	return out
}

func (t TagID) addID(container Container) {
	if _, ok := t.iDs[container.ImageID]; !ok {
		t.iDs[container.ImageID] = newID(container.ImageID)
	}
	t.iDs[container.ImageID].addContainer(container)
}

// --- container digest/id ---

// ID is container image digest/ID, it comes from container status after image is pulled
type ID struct {
	Name       string
	containers map[containerKey]Container
}

type containerKey struct {
	podName       string
	podNamespace  string
	containerName string
}

func newID(name string) ID {
	return ID{Name: name, containers: map[containerKey]Container{}}
}

func (i ID) addContainer(container Container) {
	key := containerKey{
		podName:       container.Pod.Name,
		podNamespace:  container.Pod.Namespace,
		containerName: container.Name,
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
