package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

// GetReleaseFn ...
type GetReleaseFn = func(owner, repo, tag string) (Release, error)

// ResolveBinaryPathFn ...
type ResolveBinaryPathFn = func() (string, error)

// CreateUpdateCommand creates the 'config' sub command
func CreateUpdateCommand(owner, repo, version, binaryName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Long:  fmt.Sprintf(`Checks for a newer release on GitHub and updates if one is found (https://github.com/%s/%s/releases)`, owner, repo),
		Short: `Checks for a newer release on GitHub and updates if one is found`,
		Run:   RunSelfUpdateFn(owner, repo, version, binaryName),
	}

	cmd.Flags().String("tag", "", `the version tag to update to`)

	return cmd
}

// RunSelfUpdateFn runs the self update command based on the current version and binary name.
// currentVersion is used to determine whether a newer one is available
func RunSelfUpdateFn(owner, repo, currentVersion, binaryName string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		tag, _ := cmd.Flags().GetString("tag")
		if err := RunSelfUpdate(owner, repo, tag, currentVersion, binaryName, os.Executable, GetLatestRelease); err != nil {
			log.Error(err)
			log.Exit(1)
		}
	}
}

// RunSelfUpdate checks whether the github.com/owner/repo has a release that is more recent that the specified version. If one exists, tries to
// find a binary asset that matches the current platform and the provided binaryName. If one is found, it is downloaded to the path of the current
// executable.
func RunSelfUpdate(owner, repo, requestedTag, version, binaryName string, resolveBinaryPathFn ResolveBinaryPathFn, getReleaseFn GetReleaseFn) (err error) {
	var binaryPath string
	if binaryPath, err = resolveBinaryPathFn(); err != nil {
		return err
	}

	if requestedTag == version {
		log.Infof("You already run %s version %s", binaryName, requestedTag)
		return
	}

	log.Infof("Fetching release...")
	var release Release
	if release, err = getReleaseFn(owner, repo, requestedTag); err != nil {
		return err
	}

	foundTag := release.TagName()
	if requestedTag == "" {
		log.Infof("Latest tag is %s", foundTag)
		log.Infof("Current version is %s", version)

		if foundTag != "" && foundTag != version && semver.Compare(foundTag, version) > 0 {
			err = download(foundTag, binaryName, binaryPath, release)
		} else {
			log.Infof("You are already running the latest version of %s!", binaryName)
		}
	} else {
		err = download(requestedTag, binaryName, binaryPath, release)
	}

	return err
}

func download(tagName, assetName, targetPath string, release Release) (err error) {
	log.Infof("Downloading version %s...", tagName)
	var rc io.ReadCloser
	if rc, err = release.DownloadBinary(assetName); err == nil {
		var content []byte
		if content, err = ioutil.ReadAll(rc); err == nil {
			err = ioutil.WriteFile(targetPath, content, 0755)
			defer rc.Close()
		}
	}
	return
}
