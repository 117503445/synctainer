package gh

import (
	"fmt"
	"github.com/imroc/req/v3"
)

func TriggerGithubAction() error {
	req.DevMode()

	resp, err :=

		req.SetBodyJsonMarshal(map[string]interface{}{
			"ref": "master",
			"inputs": map[string]string{
				"image":    "mysql",
				"platform": "linux/amd64",
			},
		}).SetHeader("Accept", "application/vnd.github+json").SetHeader("Authorization", "Bearer ").SetHeader("X-GitHub-Api-Version", "2022-11-28").
			Post("https://api.github.com/repos/117503445/container-copier/actions/workflows/copy.yml/dispatches")
	if err != nil {
		return err
	}

	fmt.Println(resp.String())
	return nil
}
