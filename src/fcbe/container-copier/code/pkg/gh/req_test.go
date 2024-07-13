package gh_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/117503445/container-copier/src/fcbe/pkg/gh"
)

func TestTriggerGithubAction(t *testing.T) {
	ast := assert.New(t)

	err := gh.TriggerGithubAction("mysql", "linux/amd64")

	ast.NoError(err)
}
