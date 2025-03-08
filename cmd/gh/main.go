package main

import (
	"io"
	"os"
	"strings"

	"github.com/117503445/goutils"
	"github.com/117503445/goutils/gexec"
	"github.com/117503445/synctainer/pkg/cfg"
	"github.com/117503445/synctainer/pkg/convert"

	"github.com/regclient/regclient/types/ref"
	"github.com/rs/zerolog/log"
)

func main() {
	goutils.InitZeroLog()

	cfg.CfgLoad("gh")

	c := cfg.Cfg

	srcImage := ""
	if strings.Contains(c.Image, "sha256") {
		srcImage = c.Image
	} else {
		srcRef, err := ref.New(c.Image)
		if err != nil {
			log.Fatal().Err(err).Msg("ref.New")
		}
		log.Debug().Str("srcRef", srcRef.CommonName()).Msg("")

		output, err := gexec.Run(
			gexec.Commands([]string{
				"regctl", "image", "digest", c.Image, "--platform", c.Platform,
			}),
			&gexec.RunCfg{
				Writers: []io.Writer{os.Stdout},
			},
		)
		if err != nil {
			log.Fatal().Err(err).Msg("Exec")
		}

		digest := strings.TrimSpace(output)
		if !strings.HasPrefix(digest, "sha256:") {
			log.Fatal().Str("digest", digest).Msg("digest not start with sha256:")
		}
		srcRef.Digest = digest
		srcImage = srcRef.CommonName()
	}

	newImage, err := convert.ConvertToNewImage(c.Image, c.Platform)
	if err != nil {
		log.Fatal().Err(err).Msg("ConvertToNewImage")
	}
	log.Info().Str("srcImage", srcImage).Str("newImage", newImage).Msg("ConvertToNewImage")

	cmds := []string{
		"regctl", "registry", "login", c.Registry, "--user", c.Username, "--pass", c.Password,
	}
	cmdStr := strings.Join(cmds[:len(cmds)-1], " ")+" ***"

	log.Info().Str("cmd", cmdStr).Msg("Executing")
	_, err = gexec.Run(
		gexec.Commands(cmds),
		&gexec.RunCfg{
			DisableLog: true,
			Writers:    []io.Writer{os.Stdout},
		},
	)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Info().Str("cmd", cmdStr).Msg("Executed")

	_, err = gexec.Run(
		gexec.Commands([]string{
			"regctl", "image", "copy", srcImage, newImage, "--verbosity", "trace",
		}),
		&gexec.RunCfg{
			Writers: []io.Writer{os.Stdout},
		},
	)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
