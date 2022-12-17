package api

import (
	"fmt"
	"strings"
)

const (
	tagSeparator    = ":"
	digestSeparator = "@"
)

type ImageName struct {
	Registry   string
	Repository string
	Tag        string
	Digest     string
	Error      string
}

func (i ImageName) TagDigest() string {
	if i.Tag != "" {
		return i.Tag
	}
	return i.Digest
}

func ParseImageName(image string) ImageName {

	if strings.HasPrefix(image, ":") || strings.HasSuffix(image, ":") {
		return ImageName{Error: fmt.Sprintf("image name %s is not valid", image)}
	}

	var img ImageName
	var imageWithoutVersion string
	imageWithoutVersion, img.Tag, img.Digest = getImageTagAndDigest(image)

	// get repository if it exists
	imageParts := strings.Split(imageWithoutVersion, "/")
	if len(imageParts) > 1 &&
		(strings.HasPrefix(imageParts[0], "localhost") || strings.Contains(imageParts[0], ".")) {
		// docker.io is default registry and is not included in repository names, if we included this
		// then it is hard to get image id (image id has to be searched by repository name)
		if imageParts[0] != "docker.io" {
			img.Registry = imageParts[0]
		}
		img.Repository = strings.Join(imageParts[1:], "/")
		return img
	}

	img.Repository = imageWithoutVersion
	return img
}

func getImageTagAndDigest(image string) (string, string, string) {
	if strings.Contains(image, digestSeparator) {
		imageWithoutSuffix, suffix := splitImageAndSuffix(image, digestSeparator)
		return imageWithoutSuffix, "", suffix
	}
	if strings.Contains(image, tagSeparator) {
		imageWithoutSuffix, suffix := splitImageAndSuffix(image, tagSeparator)
		return imageWithoutSuffix, suffix, ""
	}
	return image, "", ""
}

func splitImageAndSuffix(image, separator string) (string, string) {
	imageParts := strings.Split(image, separator)
	if len(imageParts) == 1 {
		return image, ""
	}
	if len(imageParts) == 2 {
		return imageParts[0], imageParts[1]
	}

	lastIndex := len(imageParts) - 1
	suffix := imageParts[lastIndex]
	imageWithoutSuffix := strings.Join(imageParts[:lastIndex], separator)
	return imageWithoutSuffix, suffix
}
