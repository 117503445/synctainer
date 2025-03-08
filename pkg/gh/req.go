package gh

import (
	"errors"

	"github.com/imroc/req/v3"
)

func TriggerGithubAction(image string, platform string, githubToken string) error {
	if platform == "" {
		platform = "linux/amd64"
	}

	req.DevMode()

	resp, err :=
		req.SetBodyJsonMarshal(map[string]any{
			"ref": "master",
			"inputs": map[string]string{
				"image":    image,
				"platform": platform,
			},
		}).SetHeader("Accept", "application/vnd.github+json").SetHeader("Authorization", "Bearer "+githubToken).SetHeader("X-GitHub-Api-Version", "2022-11-28").
			Post("https://api.github.com/repos/117503445/synctainer/actions/workflows/copy.yml/dispatches")
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return errors.New("failed to trigger github action, resp: " + resp.String())
	}

	return nil
}
