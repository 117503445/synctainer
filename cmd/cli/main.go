package main

import (
	"context"
	"net/http"
	"time"

	"github.com/117503445/goutils"
	"github.com/117503445/synctainer/pkg/rpc"
	"github.com/rs/zerolog/log"
)

func init() {
	goutils.InitZeroLog()

}

func (c *cmdSyncImage) Run() error {
	client := rpc.NewFcProtobufClient(
		cli.BackendHost, http.DefaultClient,
	)
	postResp, err := client.PostTask(context.Background(), &rpc.ReqPostTask{
		Image:    c.Image,
		Platform: c.Platform,
		Registry: cli.RegistryHost,
		Username: cli.RegistryUser,
		Password: cli.RegistryPass,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("PostTask")
		return err
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
			return err
		}

		log.Debug().Interface("getResp", getResp).Send()
		if getResp.Digest != "" {
			log.Info().Str("hashImage", getResp.Digest).Msg("Task started")
			break
		}
	}

	return nil
}

func main() {
	log.Info().Msg("Cli")
	// 1. 输入 image 名称，调用函数计算，输出 image 链接
	// 2. 输入 compose 文件，写入 override.yaml 文件
	CliLoad()
}
