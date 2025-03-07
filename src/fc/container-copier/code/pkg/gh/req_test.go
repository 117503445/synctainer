package gh_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/117503445/synctainer/src/fc/pkg/gh"
)

func TestTriggerGithubAction(t *testing.T) {
	ast := assert.New(t)

	err := gh.TriggerGithubAction("mysql", "linux/amd64")

	ast.NoError(err)
}
