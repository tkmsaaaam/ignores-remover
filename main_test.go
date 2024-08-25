package main

import (
	"io/fs"
	"testing"
	"time"
)

var isDir bool

type Impl struct{}

func (impl Impl) IsDir() bool {
	return isDir
}
func (impl Impl) Name() string {
	return "/path/to"
}
func (impl Impl) Size() int64 {
	return 0
}
func (impl Impl) Mode() fs.FileMode {
	return 0
}
func (impl Impl) ModTime() time.Time {
	return time.Now()
}
func (impl Impl) Sys() any {
	return nil
}

func TestMakeRequest(t *testing.T) {
	tests := []struct {
		name    string
		request string
		isDir   bool
		want    Request
	}{
		{
			name:    "NoArg",
			request: "",
			isDir:   true,
			want:    Request{path: "/path/to/", file: ".gitignore"},
		},
		{
			name:    "ArgIsPath",
			request: "/path/to",
			isDir:   true,
			want:    Request{path: "/path/to/", file: ".gitignore"},
		},
		{
			name:    "ArgIsPathSuffixTrue",
			request: "/path/to/",
			isDir:   true,
			want:    Request{path: "/path/to/", file: ".gitignore"},
		},
		{
			name:    "ArgIsFile",
			request: ".dockerignore",
			isDir:   false,
			want:    Request{path: "/path/to/", file: ".dockerignore"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			originalOsGetwd := osGetwd
			defer func() { osGetwd = originalOsGetwd }()
			osGetwd = func() (string, error) {
				return "/path/to", nil
			}

			originalOsStat := osStat
			defer func() { osStat = originalOsStat }()
			osStat = func(name string) (fs.FileInfo, error) {
				isDir = tt.isDir
				return Impl{}, nil
			}

			actual := makeRequest(tt.request)

			if actual.path != tt.want.path {
				t.Errorf("makeRequest() path = %v, want %v", actual.path, tt.want.path)
			}
			if actual.file != tt.want.file {
				t.Errorf("makeRequest() file = %v, want %v", actual.path, tt.want.path)
			}
		})
	}
}
