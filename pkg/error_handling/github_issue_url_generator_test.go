package errorhandling

import (
	"fmt"
	"testing"

	"net/url"

	"github.com/sha1n/clib/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestGenerateGitHubCreateNewIssueURL(t *testing.T) {
	owner, repo, title, description := test.RandomString(), test.RandomString(), test.RandomString(), test.RandomString()

	generatedURL := GenerateGitHubCreateNewIssueURL(owner, repo, title, description)

	url, err := url.Parse(generatedURL)
	assert.NoError(t, err)

	assert.Equal(t, "https", url.Scheme)
	assert.Equal(t, "github.com", url.Host)
	assert.Equal(t, expectedPath(owner, repo), url.Path)
	assert.Equal(t, title, url.Query().Get("title"))
	assert.Equal(t, generateIssueDescription(title, description), url.Query().Get("body"))
}

func expectedPath(owner, repo string) string {
	return fmt.Sprintf("/%s/%s/issues/new", owner, repo)
}
