package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/117503445/goutils"
	"github.com/117503445/goutils/gexec"
	"github.com/117503445/synctainer/pkg/convert"
	"github.com/117503445/synctainer/pkg/rpc"

	"github.com/regclient/regclient/types/ref"
	"github.com/rs/zerolog/log"
)

func main() {
	goutils.InitZeroLog()
	CfgLoad()

	enableFc := cfg.FcCallback != ""

	newRef, err := ref.New(cfg.TargetImage)
	if err != nil {
		log.Fatal().Err(err).Msg("ref.New")
	}

	var client rpc.Fc
	if enableFc {
		client = rpc.NewFcProtobufClient(
			cfg.FcCallback, http.DefaultClient,
		)

		t := &rpc.ReqPatchTask{
			Id:                cfg.TaskId,
			GithubActionRunId: cfg.RunId,
		}
		log.Info().Str("githubActionRunId", cfg.RunId).Msg("PatchTask")
		_, err := client.PatchTask(context.Background(), t)
		if err != nil {
			log.Warn().Err(err).Interface("PatchTask", t).Send()
		}
	}

	output, err := gexec.Run(
		gexec.Commands([]string{
			"regctl", "image", "digest", cfg.Image, "--platform", cfg.Platform,
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

	srcRef, err := ref.New(cfg.Image)
	if err != nil {
		log.Fatal().Err(err).Msg("ref.New")
	}
	log.Debug().Str("srcRef", srcRef.CommonName()).Msg("")
	srcRef.Digest = digest
	srcImage := srcRef.CommonName()

	newTagImage, err := convert.ConvertToNewImage(cfg.Image, cfg.Platform, cfg.TargetImage)
	if err != nil {
		log.Fatal().Err(err).Msg("ConvertToNewImage")
	}

	newHashImage := strings.Split(newTagImage, ":")[0] + "@" + digest
	log.Info().Str("srcImage", srcImage).Str("newTagImage", newTagImage).Str("newHashImage", newHashImage).Msg("ConvertToNewImage")

	if enableFc {
		t := &rpc.ReqPatchTask{
			Id:     cfg.TaskId,
			Digest: newHashImage,
		}
		log.Info().Str("digest", digest).Msg("PatchTask")
		_, err := client.PatchTask(context.Background(), t)
		if err != nil {
			log.Warn().Interface("PatchTask", t).Err(err).Send()
		}
	}

	cmds := []string{
		"regctl", "registry", "login", newRef.Registry, "--user", cfg.Username, "--pass", cfg.Password,
	}
	cmdStr := strings.Join(cmds[:len(cmds)-1], " ") + " ***"

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
			"regctl", "image", "copy", srcImage, newTagImage, "--verbosity", "info",
		}),
		&gexec.RunCfg{
			Writers: []io.Writer{os.Stdout},
		},
	)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
