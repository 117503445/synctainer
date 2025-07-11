package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/117503445/goutils"
	"github.com/117503445/synctainer/pkg/convert"
	"github.com/117503445/synctainer/pkg/gh"
	"github.com/117503445/synctainer/pkg/ots"
	"github.com/117503445/synctainer/pkg/rpc"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
)

type server struct {
	rpc.Fc
	tm          *ots.TableManager
	githubToken string
	fcCallback  string
}

func newServer(githubToken string, tm *ots.TableManager, fcCallback string) *server {
	if githubToken == "" {
		log.Fatal().Msg("github token is empty")
	}
	return &server{
		githubToken: githubToken,
		tm:          tm,
		fcCallback:  fcCallback,
	}
}

func (s *server) PostTask(ctx context.Context, req *rpc.ReqPostTask) (*rpc.RespPostTask, error) {
	tagImage, err := convert.ConvertToNewImage(req.Image, req.Platform, req.TargetImage)
	if err != nil {
		return nil, err
	}

	id := goutils.UUID7()

	err = s.tm.PutRow(id, map[string]any{})
	if err != nil {
		return nil, err
	}

	err = gh.TriggerGithubAction(id, req.Image, req.Platform, req.TargetImage, req.Username, req.Password, s.githubToken, s.fcCallback)
	if err != nil {
		return nil, err
	}

	return &rpc.RespPostTask{
		Id:       id,
		TagImage: tagImage,
	}, nil
}

func (s *server) GetTask(ctx context.Context, req *rpc.ReqGetTask) (*rpc.RespGetTask, error) {
	row, err := s.tm.GetRow(req.Id)
	if err != nil {
		log.Warn().Err(err).Msg("GetRow")
		return nil, err
	}

	runId := ots.MapMustGetString(row, "github_action_run_id")
	githubActionUrl := ""
	if runId != "" {
		githubActionUrl = fmt.Sprintf("https://github.com/117503445/synctainer/actions/runs/%v", runId)
	}

	return &rpc.RespGetTask{
		Digest:          ots.MapMustGetString(row, "digest"),
		GithubActionUrl: githubActionUrl,
	}, nil
}

func (s *server) PatchTask(ctx context.Context, req *rpc.ReqPatchTask) (*rpc.RespPatchTask, error) {
	m := map[string]any{}
	if req.Digest != "" {
		m["digest"] = req.Digest
	}
	if req.GithubActionRunId != "" {
		m["github_action_run_id"] = req.GithubActionRunId
	}
	log.Debug().Interface("m", m).Msg("PatchTask")
	err := s.tm.UpdateRow(req.Id, m)
	if err != nil {
		log.Warn().Err(err).Msg("UpdateRow")
		return nil, err
	} else {
		return &rpc.RespPatchTask{}, nil
	}
}

func main() {
	// goutils.InitZeroLog(goutils.WithNoColor{})
	goutils.InitZeroLog()
	cfgLoad()

	log.Info().Msg("Starting server...")
	tm, err := ots.NewTableManager(cfg.TablestoreEndpoint, cfg.TablestoreName, cfg.TablestoreAk, cfg.TablestoreSk)
	if err != nil {
		log.Fatal().Err(err).Msg("NewTableManager")
	}

	s := newServer(cfg.GithubToken, tm, cfg.FcCallback)
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
