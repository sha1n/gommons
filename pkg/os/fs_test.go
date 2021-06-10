package os

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandPath(t *testing.T) {
	userHomeDir, err := os.UserHomeDir()
	assert.NoError(t, err)
	aDirName := "some-directory"

	tests := []struct {
		name             string
		path             string
		wantExpandedPath string
		wantErr          bool
	}{
		{name: "no home prefix example", path: os.TempDir(), wantExpandedPath: os.TempDir(), wantErr: false},
		{name: "home prefix example", path: filepath.Join("~", aDirName), wantExpandedPath: filepath.Join(userHomeDir, aDirName), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExpandedPath, err := ExpandUserPath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExpandPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotExpandedPath != tt.wantExpandedPath {
				t.Errorf("ExpandPath() = %v, want %v", gotExpandedPath, tt.wantExpandedPath)
			}
		})
	}
}

func Test_hasHomeDirPrefix(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{name: "windows userprofile lower-case example", path: "%userprofile%\\my\\directory", want: true},
		{name: "windows userprofile upper-case example", path: "%USERPROFILE%\\my\\directory", want: true},
		{name: "windows homepath lower-case example", path: "%homepath%\\my\\directory", want: true},
		{name: "windows homepath upper-case example", path: "%HOMEPATH%\\my\\directory", want: true},
		{name: "unix home upper-case example", path: "$HOME/my/directory", want: true},
		{name: "unix home lower-case example", path: "$home/my/directory", want: false},
		{name: "unix ~ upper-case example", path: "~/my/directory", want: true},
		{name: "unix no home prefix example", path: "/my/directory", want: false},
		{name: "windows no home prefix example", path: "D:\\my\\directory", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasHomeDirPrefix(tt.path); got != tt.want {
				t.Errorf("hasHomeDirPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
