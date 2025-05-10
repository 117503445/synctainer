package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/117503445/goutils"
	"github.com/117503445/synctainer/pkg/rpc"
	"github.com/rs/zerolog/log"
)

func init() {
	goutils.InitZeroLog()
}

func syncImage(client rpc.Fc, image string, platform string) string {
	postResp, err := client.PostTask(context.Background(), &rpc.ReqPostTask{
		Image:    image,
		Platform: platform,
		TargetImage: cli.TargetImage,
		Username: cli.RegistryUser,
		Password: cli.RegistryPass,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("PostTask")
	}

	log.Debug().Interface("postResp", postResp).Msg("PostTask")

	taskId := postResp.Id

	delays := []int{10, 10, 10, 10, 10, 10, 60, 60, 60, 60}
	for _, delay := range delays {
		log.Info().Msgf("Waiting %d seconds", delay)
		time.Sleep(time.Duration(delay) * time.Second)

		getResp, err := client.GetTask(context.Background(), &rpc.ReqGetTask{
			Id: taskId,
		})
		if err != nil {
			log.Fatal().Err(err).Msg("GetTask")
		}

		log.Debug().Interface("getResp", getResp).Send()
		if getResp.Digest != "" {
			log.Info().Str("hashImage", getResp.Digest).Msg("Task started")
			return getResp.Digest
		}
	}

	return ""
}

func (c *cmdSyncImage) Run() error {
	client := rpc.NewFcProtobufClient(
		cli.BackendHost, http.DefaultClient,
	)

	syncImage(client, c.Image, c.Platform)

	return nil
}

func (c *cmdSyncCompose) Run() error {
	if !goutils.FileExists(c.ComposePath) {
		log.Fatal().Str("ComposePath", c.ComposePath).Msg("Compose file not found")
	}

	readServiceImage := func(file string) map[string]string {
		serviceImage := make(map[string]string)
		composeMap := make(map[string]any)
		err := goutils.ReadYaml(file, &composeMap)
		if err != nil {
			log.Fatal().Err(err).Msg("ReadYaml")
		}
		for serviceName, service := range composeMap["services"].(map[string]any) {
			if service.(map[string]any)["image"] != nil {
				serviceImage[serviceName] = service.(map[string]any)["image"].(string)
			}
		}
		return serviceImage
	}

	composeServiceImage := readServiceImage(c.ComposePath)
	overrideMap := make(map[string]any)
	overrideServiceImage := make(map[string]string)
	if goutils.FileExists(c.OverridePath) {
		overrideServiceImage = readServiceImage(c.OverridePath)
		err := goutils.ReadYaml(c.OverridePath, &overrideMap)
		if err != nil {
			log.Fatal().Err(err).Msg("ReadYaml")
		}
	}

	log.Debug().
		Interface("composeServiceImage", composeServiceImage).
		Interface("overrideServiceImage", overrideServiceImage).Send()

	client := rpc.NewFcProtobufClient(
		cli.BackendHost, http.DefaultClient,
	)

	for service, image := range composeServiceImage {
		if _, ok := overrideServiceImage[service]; ok {
			continue
		}

		hashImage := syncImage(client, image, c.Platform)
		if overrideMap["services"] == nil {
			overrideMap["services"] = map[string]any{}
		}

		overrideMap["services"].(map[string]any)[service] = map[string]any{
			"image": hashImage,
		}
	}

	if goutils.FileExists(c.OverridePath) {
		src := c.OverridePath
		dest := fmt.Sprintf("%s.%s.bak", c.OverridePath, goutils.TimeStrMilliSec())
		err := goutils.CopyFile(src, dest)
		if err != nil {
			log.Fatal().Err(err).
				Str("src", src).
				Str("dest", dest).
				Msg("backup")
		}
	}

	err := goutils.WriteYaml(c.OverridePath, overrideMap)
	if err != nil {
		log.Fatal().Err(err).Msg("WriteYaml")
	}

	return nil
}

func main() {
	log.Info().Msg("Cli")
	// 1. 输入 image 名称，调用函数计算，输出 image 链接
	// 2. 输入 compose 文件，写入 override.yaml 文件
	CliLoad()
}
