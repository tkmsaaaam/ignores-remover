package main

import (
	"errors"
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
	type GetWd struct {
		path string
		err  error
	}
	type Stat struct {
		isDir bool
		err   error
	}

	tests := []struct {
		name    string
		request string
		stat    Stat
		getWd   GetWd
		want    *Request
	}{
		{
			name:    "NoArg",
			request: "",
			stat:    Stat{isDir: true, err: nil},
			getWd:   GetWd{path: "/path/to", err: nil},
			want:    &Request{path: "/path/to/", file: ".gitignore"},
		},
		{
			name:    "NoArgError",
			request: "",
			stat:    Stat{isDir: true, err: nil},
			getWd:   GetWd{path: "", err: errors.New("")},
			want:    nil,
		},
		{
			name:    "ArgIsPath",
			request: "/path/to",
			stat:    Stat{isDir: true, err: nil},
			getWd:   GetWd{path: "/path/to", err: nil},
			want:    &Request{path: "/path/to/", file: ".gitignore"},
		},
		{
			name:    "ArgIsPathSuffixTrue",
			request: "/path/to/",
			stat:    Stat{isDir: true, err: nil},
			getWd:   GetWd{path: "/path/to", err: nil},
			want:    &Request{path: "/path/to/", file: ".gitignore"},
		},
		{
			name:    "ArgIsFullPathToFile",
			request: "/path/to/.dockerignore",
			stat:    Stat{isDir: false, err: nil},
			getWd:   GetWd{path: "/path/to", err: nil},
			want:    &Request{path: "/path/to/", file: ".dockerignore"},
		},
		{
			name:    "ArgIsFile",
			request: ".dockerignore",
			stat:    Stat{isDir: false, err: nil},
			getWd:   GetWd{path: "/path/to", err: nil},
			want:    &Request{path: "/path/to/", file: ".dockerignore"},
		},
		{
			name:    "ArgIsFileError",
			request: ".dockerignore",
			stat:    Stat{isDir: false, err: nil},
			getWd:   GetWd{path: "", err: errors.New("")},
			want:    nil,
		},
		{
			name:    "statIsError",
			request: "/path/to",
			stat:    Stat{isDir: false, err: errors.New("")},
			getWd:   GetWd{path: "/path/to/", err: nil},
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			originalOsGetwd := osGetwd
			defer func() { osGetwd = originalOsGetwd }()
			osGetwd = func() (string, error) {
				return tt.getWd.path, tt.getWd.err
			}

			originalOsStat := osStat
			defer func() { osStat = originalOsStat }()
			osStat = func(name string) (fs.FileInfo, error) {
				isDir = tt.stat.isDir
				return Impl{}, tt.stat.err
			}

			actual := makeRequest(tt.request)

			if tt.want != nil {
				if actual.path != tt.want.path {
					t.Errorf("makeRequest() path = %v, want %v", actual.path, tt.want.path)
				}
				if actual.file != tt.want.file {
					t.Errorf("makeRequest() file = %v, want %v", actual.path, tt.want.path)
				}
			} else {
				if actual != nil {
					t.Errorf("makeRequest() is not nil %v", actual)
				}
			}

		})
	}
}
