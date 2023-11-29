package api

import (
	v1 "k8s.io/api/core/v1"
	"time"
)

// --- nodes ---

func NewNodes(nodes []v1.Node) map[string]Node {
	out := make(map[string]Node)
	for _, node := range nodes {
		out[node.Name] = NewNode(node)
	}
	return out
}

// --- node ---

type Node struct {
	Name       string
	Created    time.Time
	NodeImages NodeImages
}

func NewNode(node v1.Node) Node {
	return Node{
		Name:       node.Name,
		Created:    node.CreationTimestamp.Time,
		NodeImages: newNodeImages(node.Status.Images),
	}
}

// --- node images ---

type NodeImages []nodeImage

func newNodeImages(images []v1.ContainerImage) NodeImages {
	var out NodeImages
	for _, image := range images {
		out = append(out, newNodeImage(image))
	}
	return out
}

// GetSizeBytes return size in bytes if the image is found on this node, 0 is returned otherwise
func (n NodeImages) GetSizeBytes(imageName ImageName) int64 {
	for _, image := range n {
		if size := image.getSizeBytes(imageName); size != 0 {
			return size
		}
	}
	return 0
}

// --- node image ---

type nodeImage struct {
	names     []ImageName
	sizeBytes int64
}

func newNodeImage(image v1.ContainerImage) nodeImage {
	var names []ImageName
	for _, name := range image.Names {
		names = append(names, ParseImageName(name))
	}
	return nodeImage{
		names:     names,
		sizeBytes: image.SizeBytes,
	}
}

// sizeBytes return size in bytes if the image is found, 0 is returned otherwise
func (n nodeImage) getSizeBytes(imageName ImageName) int64 {
	for _, name := range n.names {
		if name.Equals(imageName) {
			return n.sizeBytes
		}
	}
	return 0
}
