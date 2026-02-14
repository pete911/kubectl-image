package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseImageName(t *testing.T) {

	t.Run("image without registry", func(t *testing.T) {
		image := ParseImageName("memcached:1.5")

		assert.Equal(t, "", image.Error)
		assert.Equal(t, "", image.Registry)
		assert.Equal(t, "memcached", image.Repository)
		assert.Equal(t, "1.5", image.Tag)
		assert.Equal(t, "", image.ID)
	})

	t.Run("when image has docker.io registry then registry is removed", func(t *testing.T) {
		image := ParseImageName("docker.io/memcached:1.5")

		assert.Equal(t, "", image.Error)
		assert.Equal(t, "", image.Registry)
		assert.Equal(t, "memcached", image.Repository)
		assert.Equal(t, "1.5", image.Tag)
		assert.Equal(t, "", image.ID)
	})

	t.Run("image with registry", func(t *testing.T) {
		image := ParseImageName("quay.io/jacksontj/promxy:v0.0.58")

		assert.Equal(t, "", image.Error)
		assert.Equal(t, "quay.io", image.Registry)
		assert.Equal(t, "jacksontj/promxy", image.Repository)
		assert.Equal(t, "v0.0.58", image.Tag)
		assert.Equal(t, "", image.ID)
	})

	t.Run("image with tag and digest/ID", func(t *testing.T) {
		image := ParseImageName("public.ecr.aws/karpenter/controller:v0.31.1@sha256:be5da7cb9181c2969280e869c465927de278e17d56762bfa5eeb4b21e5c9762e")

		assert.Equal(t, "", image.Error)
		assert.Equal(t, "public.ecr.aws", image.Registry)
		assert.Equal(t, "karpenter/controller", image.Repository)
		assert.Equal(t, "v0.31.1", image.Tag)
		assert.Equal(t, "sha256:be5da7cb9181c2969280e869c465927de278e17d56762bfa5eeb4b21e5c9762e", image.ID)
	})

	t.Run("image with digest/ID", func(t *testing.T) {
		image := ParseImageName("gcr.io/google-containers/pause-amd64@sha256:4a1c4b21597c1b4415bdbecb28a3296c6b5e23ca4f9feeb599860a1dac6a0108")

		assert.Equal(t, "", image.Error)
		assert.Equal(t, "gcr.io", image.Registry)
		assert.Equal(t, "google-containers/pause-amd64", image.Repository)
		assert.Equal(t, "", image.Tag)
		assert.Equal(t, "sha256:4a1c4b21597c1b4415bdbecb28a3296c6b5e23ca4f9feeb599860a1dac6a0108", image.ID)
	})

	t.Run("image ends with : is invalid", func(t *testing.T) {
		image := ParseImageName("quay.io/jacksontj/promxy:")
		assert.NotEmpty(t, image.Error)
	})

	t.Run("image starts with : is invalid", func(t *testing.T) {
		image := ParseImageName(":latest")
		assert.NotEmpty(t, image.Error)
	})

	t.Run("image name is : is invalid", func(t *testing.T) {
		image := ParseImageName(":")
		assert.NotEmpty(t, image.Error)
	})
}
