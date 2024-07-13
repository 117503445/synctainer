package convert

import (
	"errors"
	"strings"

	dockerparser "github.com/novln/docker-parser"
)

const (
	NEW_REGISTRY = "registry.cn-hangzhou.aliyuncs.com"
	NEW_USERNAME = "117503445-mirror"
)

func ConvertToNewImage(image string) (string, error) {
	reference, err := dockerparser.Parse(image)
	if err != nil {
		return "", err
	}

	splits := strings.Split(reference.ShortName(), "/")
	if len(splits) <= 1 {
		return "", errors.New("image with len(shortName.split(\"/\")) <= 1")
	}

	// if len(splits) != 2 {
	// 	log.Warn().Str("image", image).Msg("image with len(shortName.split(\"/\")) > 2")
	// }

	username := splits[0]
	name := strings.Join(splits[1:], "/")

	tag := reference.Tag()
	var suffix string
	if strings.HasPrefix(tag, "sha256:") {
		suffix = "@" + tag
	} else {
		suffix = ":" + tag
	}

	newImage := NEW_REGISTRY + "/" + NEW_USERNAME + "/" + reference.Registry() + "." + username + "." + name + suffix

	return newImage, nil
}
