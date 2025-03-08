package gh_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/117503445/synctainer/pkg/gh"
)

func TestTriggerGithubAction(t *testing.T) {
	ast := assert.New(t)

	err := gh.TriggerGithubAction("", "", "", "", "mysql", "linux/amd64", "TODO")

	ast.NoError(err)
}
