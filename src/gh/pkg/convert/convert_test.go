package convert_test

import (
	"testing"

	"github.com/117503445/container-copier/src/gh/pkg/convert"

	dockerparser "github.com/novln/docker-parser"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func referencePrint(reference *dockerparser.Reference) {
	log.Debug().Str("name", reference.Name()).Str("tag", reference.Tag()).Str("registry", reference.Registry()).Str("shortName", reference.ShortName()).Str("tag", reference.Tag()).Msg("")
}

func TestImageParse(t *testing.T) {
	ast := assert.New(t)
	reference, err := dockerparser.Parse("mysql")
	ast.NoError(err)
	referencePrint(reference)

	reference, err = dockerparser.Parse("ubuntu:18.04@sha256:98706f0f213dbd440021993a82d2f70451a73698315370ae8615cc468ac06624")
	ast.NoError(err)
	referencePrint(reference)

}

func TestConvertToNewImage(t *testing.T) {
	ast := assert.New(t)

	cases := []struct {
		image    string
		expected string
	}{
		{
			image:    "mysql",
			expected: "registry.cn-hangzhou.aliyuncs.com/117503445-mirror/sync:docker.io.library.mysql.latest",
		}, {
			image:    "mysql:8.0",
			expected: "registry.cn-hangzhou.aliyuncs.com/117503445-mirror/sync:docker.io.library.mysql.8.0",
		},
		{
			image:    "ghcr.io/devcontainers/features/anaconda:1",
			expected: "registry.cn-hangzhou.aliyuncs.com/117503445-mirror/sync:ghcr.io.devcontainers.features.anaconda.1",
		},
		{
			image:    "ubuntu:18.04@sha256:98706f0f213dbd440021993a82d2f70451a73698315370ae8615cc468ac06624",
			expected: "registry.cn-hangzhou.aliyuncs.com/117503445-mirror/sync:docker.io.library.ubuntu.sha256.98706f0f213dbd440021993a82d2f70451a73698315370ae8615cc468ac06624",
		},
	}
	for _, c := range cases {
		newImage, err := convert.ConvertToNewImage(c.image)
		ast.NoError(err)
		ast.Equal(c.expected, newImage)
	}
}
