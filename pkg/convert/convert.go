package convert

import (
	// "errors"
	"strings"

	dockerparser "github.com/novln/docker-parser"
	// "github.com/regclient/regclient/types/ref"
)

func ConvertToNewImage(image string, platform string, targetImage string) (string, error) {
	image = strings.TrimSpace(image)
	targetImage = strings.TrimSpace(targetImage)

	reference, err := dockerparser.Parse(image)
	if err != nil {
		return "", err
	}

	newTag := platform + "." +
		reference.Registry() + "." +
		reference.Name()

	chars := []string{"/", ":", "@"}
	for _, char := range chars {
		newTag = strings.ReplaceAll(newTag, char, ".")
	}

	var targetReference *dockerparser.Reference
	if targetImage != "" {
		targetReference, err = dockerparser.Parse(targetImage)
		if err != nil {
			return "", err
		}
	} else {
		// Use defaults if targetImage is empty
		defaultRegistry := "registry.cn-hangzhou.aliyuncs.com"
		defaultRepo := "117503445-mirror/sync"
		targetReference, _ = dockerparser.Parse(defaultRegistry + "/" + defaultRepo)
	}

	// Avoid prepending 'docker.io' for non-namespaced repos like my-registry/my-repo
	registry := targetReference.Registry()
	if registry == "docker.io" {
		// Check if the input targetImage contains a slash but no explicit registry
		if strings.Contains(targetImage, "/") && !strings.Contains(targetImage, ".") {
			registry = ""
		}
	} else {
		registry = registry + "/"
	}

	newImage := registry + targetReference.ShortName() + ":" + newTag
	newImage = strings.TrimPrefix(newImage, "//") // remove double slashes if any

	return newImage, nil
}
