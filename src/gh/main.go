package main

import (
	"os"

	"github.com/117503445/synctainer/src/gh/pkg/convert"
	"github.com/117503445/goutils"
	
	"github.com/rs/zerolog/log"
)

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

	newImage, err := convert.ConvertToNewImage(image)
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
