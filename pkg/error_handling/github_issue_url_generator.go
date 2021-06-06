package errorhandling

import (
	"fmt"
	"net/url"
	"runtime"
)

// GenerateGitHubCreateNewIssueURL generates a GitHub issue URL for the specified owner/repo
// using the provided title and description.
// In addition to the provided description adds environment contextual information.
func GenerateGitHubCreateNewIssueURL(owner, repo, title, description string) string {
	params := url.Values{}
	params.Set("title", title)
	params.Set("body", generateIssueDescription(title, description))

	return generateBaseURL(owner, repo) + params.Encode()
}

func generateBaseURL(owner, repo string) string {
	return fmt.Sprintf("https://github.com/%s/%s/issues/new?", owner, repo)
}

func generateIssueDescription(title, description string) string {
	return fmt.Sprintf(`# %s
%s

GOOS: %s, GOARCH: %s, CPUs: %d, Go routines: %d

`, title, description, runtime.GOOS, runtime.GOARCH, runtime.NumCPU(), runtime.NumGoroutine())
}
