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

	id := goutils.UUID7()

	err := gh.TriggerGithubAction(id, req.Image, req.Platform, req.Registry, req.Username, req.Password, s.githubToken, "TODO")
	if err != nil {
		return nil, err
	}

	return &rpc.RespPostTask{
		Id: id,
	}, nil
}

func main() {
	// goutils.InitZeroLog(goutils.WithNoColor{})
	goutils.InitZeroLog()

	cfg.CfgLoad("fc")

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
