package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/117503445/goutils"
	"github.com/117503445/goutils/gexec"
	"github.com/117503445/synctainer/src/gh/pkg/convert"

	"github.com/regclient/regclient/types/ref"
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

	srcImage := ""
	if strings.Contains(image, "sha256") {
		srcImage = image
	} else {
		srcRef, err := ref.New(image)
		if err != nil {
			log.Fatal().Err(err).Msg("ref.New")
		}
		log.Debug().Str("srcRef", srcRef.CommonName()).Msg("")

		result, err := goutils.Exec(fmt.Sprintf("./regctl image digest %v --platform %v", image, platform))
		if err != nil {
			log.Fatal().Err(err).Msg("Exec")
		}
		digest := strings.TrimSpace(result.Stdout)
		if !strings.HasPrefix(digest, "sha256:") {
			log.Fatal().Str("digest", digest).Msg("digest not start with sha256:")
		}
		srcRef.Digest = digest
		srcImage = srcRef.CommonName()
	}

	newImage, err := convert.ConvertToNewImage(image, platform)
	if err != nil {
		log.Fatal().Err(err).Msg("ConvertToNewImage")
	}

	log.Info().Str("srcImage", srcImage).Str("newImage", newImage).Msg("ConvertToNewImage")

	_, err = gexec.Run(
		gexec.Commands([]string{
			"./regctl", "image", "copy", srcImage, newImage, "--verbosity", "trace",
		}),
		&gexec.RunCfg{
			Writers: []io.Writer{os.Stdout},
		},
	)

	// _, err = goutils.Exec(fmt.Sprintf("./regctl image copy %v %v", srcImage, newImage), goutils.WithDumpOutput{})
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
