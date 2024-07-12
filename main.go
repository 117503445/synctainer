package main

import (
	"errors"
	"os"
	"strings"

	"github.com/117503445/goutils"
	dockerparser "github.com/novln/docker-parser"

	"github.com/rs/zerolog/log"
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

func main() {
	goutils.InitZeroLog()

	image := os.Getenv("IMAGE")
	if image == "" {
		log.Fatal().Msg("Image is required")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		platform = "linux/amd64"
	}

	log.Info().Str("image", image).Str("platform", platform).Msg("load env")

	newImage, err := ConvertToNewImage(image)
	if err != nil {
		log.Fatal().Err(err).Msg("ConvertToNewImage")
	}

	log.Info().Str("newImage", newImage).Msg("ConvertToNewImage")

	if err := goutils.CMD("", "docker", "pull", "--platform", platform, image); err != nil {
		log.Fatal().Err(err).Msg("docker pull")
	}

	if err := goutils.CMD("", "docker", "tag", image, newImage); err != nil {
		log.Fatal().Err(err).Msg("docker tag")
	}

	if err := goutils.CMD("", "docker", "push", newImage); err != nil {
		log.Fatal().Err(err).Msg("docker push")
	}
}