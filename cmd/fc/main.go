package main

import (
	"context"
	"net/http"

	"github.com/117503445/goutils"
	"github.com/117503445/synctainer/pkg/cfg"
	"github.com/117503445/synctainer/pkg/gh"
	"github.com/117503445/synctainer/pkg/rpc"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
)

type server struct {
	rpc.Fc
	githubToken string
}

func newServer(githubToken string) *server {
	if githubToken == "" {
		log.Fatal().Msg("github token is empty")
	}
	return &server{
		githubToken: githubToken,
	}
}

func (s *server) PostTask(ctx context.Context, req *rpc.ReqPostTask) (*rpc.RespPostTask, error) {
	// newImage, err := convert.ConvertToNewImage(req.Image, req.Platform)
	// if err != nil {
	// 	return nil, err
	// }

	err := gh.TriggerGithubAction(req.Image, req.Platform, s.githubToken)
	if err != nil {
		return nil, err
	}

	// TODO: uuid7
	id := goutils.UUID4()

	return &rpc.RespPostTask{
		Id: id,
	}, nil
}

func main() {
	// goutils.InitZeroLog(goutils.WithNoColor{})
	goutils.InitZeroLog()

	cfg.CfgLoad()

	log.Info().Msg("Starting server...")
	s := newServer(cfg.Cfg.GithubToken)
	var handler http.Handler

	handler = rpc.NewFcServer(s)
	corsWrapper := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST"},
		AllowedHeaders: []string{"Content-Type"},
	})
	handler = corsWrapper.Handler(handler)

	http.Handle("/", handler)
	http.ListenAndServe(":8080", nil)
}
