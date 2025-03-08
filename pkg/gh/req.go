package gh

import (
	"errors"

	"github.com/imroc/req/v3"
)

func TriggerGithubAction(taskId string, image string, platform string, registry string, username string, password string, githubToken string, fcCallback string) error {
	if platform == "" {
		platform = "linux/amd64"
	}

	req.DevMode()

	resp, err :=
		req.SetBodyJsonMarshal(map[string]any{
			"ref": "master",
			"inputs": map[string]string{
				"taskId":     taskId,
				"fcCallback": fcCallback,
				"image":      image,
				"platform":   platform,
				"registry":   registry,
				"username":   username,
				"password":   password,
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
