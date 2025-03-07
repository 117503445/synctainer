package convert

import (
	// "errors"
	"strings"

	dockerparser "github.com/novln/docker-parser"
	// "github.com/regclient/regclient/types/ref"
)

const (
	NEW_REGISTRY  = "registry.cn-hangzhou.aliyuncs.com"
	NEW_SHORTNAME = "117503445-mirror/sync"
)

// func ConvertToSrcImage(image string, digest string) (string,error){
// 	ref ,err:= ref.New(image)
// 	if err != nil {
// 		return "", err
// 	}
// }

func ConvertToNewImage(image string, platform string) (string, error) {
	image = strings.TrimSpace(image)

	reference, err := dockerparser.Parse(image)
	if err != nil {
		return "", err
	}

	newTag := platform + "." +
		reference.Registry() + "." +
		reference.Name()

	chars := []string{"/", ":", "@"}
	for _, char := range chars {
		newTag = strings.ReplaceAll(newTag, char, ".")
	}

	newImage := NEW_REGISTRY + "/" + NEW_SHORTNAME + ":" + newTag

	return newImage, nil
}
