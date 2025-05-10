package convert_test

import (
	"testing"

	"github.com/117503445/synctainer/pkg/convert"

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
		image        string
		targetImage  string
		expected     string
		expectErr    bool
	}{
		{
			image:       "mysql",
			targetImage: "registry.cn-hangzhou.aliyuncs.com/117503445-mirror/sync",
			expected:    "registry.cn-hangzhou.aliyuncs.com/117503445-mirror/sync:linux.amd64.docker.io.library.mysql.latest",
		},
		{
			image:       "mysql:8.0",
			targetImage: "registry.cn-hangzhou.aliyuncs.com/117503445-mirror/sync",
			expected:    "registry.cn-hangzhou.aliyuncs.com/117503445-mirror/sync:linux.amd64.docker.io.library.mysql.8.0",
		},
		{
			image:       "ghcr.io/github/super-linter:v5",
			targetImage: "my-registry/my-repo",
			expected:    "my-registry/my-repo:linux.amd64.ghcr.io.github.super-linter.v5",
		},
		{
			image:       "ubuntu:18.04@sha256:98706f0f213dbd440021993a82d2f70451a73698315370ae8615cc468ac06624",
			targetImage: "custom-reg/custom-repo",
			expected:    "custom-reg/custom-repo:linux.amd64.docker.io.library.ubuntu.sha256.98706f0f213dbd440021993a82d2f70451a73698315370ae8615cc468ac06624",
		},
		{
			image:       "invalid-image",
			targetImage: "not-a-valid-target:::",
			expectErr:   true,
		},
	}

	for _, c := range cases {
		newImage, err := convert.ConvertToNewImage(c.image, "linux/amd64", c.targetImage)
		if c.expectErr {
			ast.Error(err)
		} else {
			ast.NoError(err)
			ast.Equal(c.expected, newImage)
		}
	}
}