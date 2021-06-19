package cmd

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/sha1n/gommons/pkg/test"
	"github.com/stretchr/testify/assert"
)

const (
	v1_0_0 = "v1.0.0"
	v1_0_1 = "v1.0.1"
)

func TestRunSelfUpdateWithReleaseError(t *testing.T) {
	runWith(func(resolveBinaryPathFn ResolveBinaryPathFn) {
		expectedError := errors.New(test.RandomString())
		getLatestRelease := aGetReleaseFnWith(nil, expectedError)
		owner, repo := test.RandomString(), test.RandomString()

		actualError := RunSelfUpdate(owner, repo, "", test.RandomString(), test.RandomString(), resolveBinaryPathFn, getLatestRelease)

		assert.Equal(t, expectedError, actualError)
	})
}

func TestRunSelfUpdateWithCurrentRelease(t *testing.T) {
	runWith(func(resolveBinaryPathFn ResolveBinaryPathFn) {
		currentVersion := v1_0_0
		latestRelease := &fakeRelease{tag: v1_0_0}
		getLatestRelease := aGetReleaseFnWith(latestRelease, nil)
		owner, repo := test.RandomString(), test.RandomString()

		actualError := RunSelfUpdate(owner, repo, "", currentVersion, test.RandomString(), resolveBinaryPathFn, getLatestRelease)

		assert.NoError(t, actualError)
		assert.False(t, latestRelease.downloadCalled)
	})
}

func TestRunSelfUpdateWithCurrentTaggedRelease(t *testing.T) {
	runWith(func(resolveBinaryPathFn ResolveBinaryPathFn) {
		currentVersion := v1_0_0
		requestedTag := v1_0_0
		latestRelease := &fakeRelease{tag: requestedTag}
		getLatestRelease := aGetReleaseFnWith(latestRelease, nil)
		owner, repo := test.RandomString(), test.RandomString()

		actualError := RunSelfUpdate(owner, repo, requestedTag, currentVersion, test.RandomString(), resolveBinaryPathFn, getLatestRelease)

		assert.NoError(t, actualError)
		assert.False(t, latestRelease.downloadCalled)
	})
}

func TestRunSelfUpdateWithDownloadError(t *testing.T) {
	runWith(func(resolveBinaryPathFn ResolveBinaryPathFn) {
		currentVersion := v1_0_0
		expectedError := errors.New(test.RandomString())
		latestRelease := &fakeRelease{tag: v1_0_1, downloadError: expectedError}
		getLatestRelease := aGetReleaseFnWith(latestRelease, nil)
		owner, repo := test.RandomString(), test.RandomString()

		actualError := RunSelfUpdate(owner, repo, "", currentVersion, test.RandomString(), resolveBinaryPathFn, getLatestRelease)

		assert.Error(t, actualError)
		assert.Equal(t, expectedError, actualError)
		assert.True(t, latestRelease.downloadCalled)
	})
}

func TestRunSelfUpdateWithSuccessfulDownload(t *testing.T) {
	runWith(func(resolveBinaryPathFn ResolveBinaryPathFn) {
		expectedFileContent := []byte(test.RandomString())
		currentVersion := v1_0_0
		latestRelease := &fakeRelease{tag: v1_0_1, data: expectedFileContent}
		getLatestRelease := aGetReleaseFnWith(latestRelease, nil)
		owner, repo := test.RandomString(), test.RandomString()

		actualError := RunSelfUpdate(owner, repo, "", currentVersion, test.RandomString(), resolveBinaryPathFn, getLatestRelease)

		assert.NoError(t, actualError)
		assert.True(t, latestRelease.downloadCalled)

		path, err := resolveBinaryPathFn()
		assert.NoError(t, err)

		actualFileContent, err := os.ReadFile(path)
		assert.NoError(t, err)

		assert.Equal(t, expectedFileContent, actualFileContent)
	})
}

func TestRunSelfUpdateWithTaggedSuccessfulDownload(t *testing.T) {
	runWith(func(resolveBinaryPathFn ResolveBinaryPathFn) {
		expectedFileContent := []byte(test.RandomString())
		currentVersion := v1_0_1
		requestedTag := v1_0_0
		requestedRelease := &fakeRelease{tag: requestedTag, data: expectedFileContent}
		getLatestRelease := aGetReleaseFnWith(requestedRelease, nil)
		owner, repo := test.RandomString(), test.RandomString()

		actualError := RunSelfUpdate(owner, repo, requestedTag, currentVersion, test.RandomString(), resolveBinaryPathFn, getLatestRelease)

		assert.NoError(t, actualError)
		assert.True(t, requestedRelease.downloadCalled)

		path, err := resolveBinaryPathFn()
		assert.NoError(t, err)

		actualFileContent, err := os.ReadFile(path)
		assert.NoError(t, err)

		assert.Equal(t, expectedFileContent, actualFileContent)
	})
}

func TestRunSelfUpdateWithAlreadyExistingRequestedVersion(t *testing.T) {
	runWith(func(resolveBinaryPathFn ResolveBinaryPathFn) {
		currentVersion := v1_0_1
		requestedTag := v1_0_1
		requestedRelease := &fakeRelease{tag: requestedTag, data: nil}
		getLatestRelease := aGetReleaseFnWith(requestedRelease, nil)
		owner, repo := test.RandomString(), test.RandomString()

		actualError := RunSelfUpdate(owner, repo, requestedTag, currentVersion, test.RandomString(), resolveBinaryPathFn, getLatestRelease)

		assert.NoError(t, actualError)
		assert.False(t, requestedRelease.downloadCalled)
	})
}

func aGetReleaseFnWith(r Release, e error) GetReleaseFn {
	return func(owner, repo, tag string) (Release, error) {
		return r, e
	}
}

func runWith(doTest func(resolveBinaryPathFn ResolveBinaryPathFn)) {
	resolveBinaryPathFn, cleanup := resolveExecutableFn()
	defer cleanup()

	doTest(resolveBinaryPathFn)
}

type fakeRelease struct {
	tag            string
	downloadError  error
	downloadCalled bool
	data           []byte
}

func (r *fakeRelease) TagName() string {
	return r.tag
}

func (r *fakeRelease) DownloadBinary(binaryName string) (rc io.ReadCloser, err error) {
	r.downloadCalled = true
	rc, err = nil, r.downloadError

	if r.downloadError == nil {
		rc = ioutil.NopCloser(bytes.NewReader(r.data))
	}

	return rc, err
}

func resolveExecutableFn() (ResolveBinaryPathFn, func()) {
	f, _ := ioutil.TempFile("", "fake_binary")

	fn := func() (string, error) {
		return f.Name(), nil
	}

	return fn, func() { os.Remove(f.Name()) }
}
