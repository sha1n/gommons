package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"runtime"

	"github.com/google/go-github/v35/github"
)

// Release a GitHub realease facade
type Release interface {
	TagName() string
	DownloadBinary(string) (io.ReadCloser, error)
}

type rel struct {
	owner    string
	repo     string
	delegate *github.RepositoryRelease
}

// GetLatestRelease returns the latest github non-draft release of this program.
func GetLatestRelease(owner, repo, tag string) (release Release, err error) {
	ctx := context.Background()
	client := github.NewClient(nil)

	var response *github.Response
	var ghRelease *github.RepositoryRelease
	if tag == "" {
		ghRelease, response, err = client.Repositories.GetLatestRelease(ctx, owner, repo)
	} else {
		ghRelease, response, err = client.Repositories.GetReleaseByTag(ctx, owner, repo, tag)
	}

	switch response.StatusCode {
	case 404:
		err = errors.New("no matching release could be found")

	case 200:
		release = &rel{
			owner:    owner,
			repo:     repo,
			delegate: ghRelease,
		}
	}

	return
}

func (r *rel) TagName() string {
	return *r.delegate.TagName
}

func (r *rel) DownloadBinary(binaryName string) (rc io.ReadCloser, err error) {
	ctx := context.Background()
	client := github.NewClient(nil)
	var assetID int64

	if assetID, err = findCompatibleAssetID(binaryName, r.delegate); err == nil {
		rc, _, err = client.Repositories.DownloadReleaseAsset(ctx, r.owner, r.repo, assetID, http.DefaultClient)
	}

	return rc, err
}

func findCompatibleAssetID(binaryName string, release *github.RepositoryRelease) (int64, error) {
	requiredAssetName := getRequiredAssetName(binaryName)
	slog.Debug(fmt.Sprintf("Required asset name is %s. Looking for matching assets in latest release.", requiredAssetName))
	for _, asset := range (*release).Assets {
		if *asset.Name == requiredAssetName {
			slog.Debug(fmt.Sprintf("Found asset ID = %d", *asset.ID))
			slog.Debug(fmt.Sprintf("Found asset Name = %s", *asset.Name))
			return *asset.ID, nil
		}
	}
	return 0, fmt.Errorf("unable to find a compatible asset in the latest release (required=%s)", requiredAssetName)
}

func getRequiredAssetName(binaryName string) string {
	assertName := fmt.Sprintf("%s-%s-%s", binaryName, runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		assertName += ".exe"
	}

	slog.Debug(fmt.Sprintf("Required asset name is: %s", assertName))

	return assertName
}
