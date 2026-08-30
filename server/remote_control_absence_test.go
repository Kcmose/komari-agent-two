package server

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRemoteFileControlIsAbsentFromSource(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate source tree")
	}
	root := filepath.Dir(filepath.Dir(thisFile))

	for _, name := range []string{
		"file_stream.go",
		"file_stream_test.go",
		"files.go",
		"files_test.go",
		"files_unix.go",
		"files_windows.go",
		"files_windows_test.go",
	} {
		path := filepath.Join(root, "server", name)
		if _, err := os.Stat(path); err == nil {
			t.Errorf("remote file-control source still exists: %s", path)
		} else if !os.IsNotExist(err) {
			t.Fatalf("stat %s: %v", path, err)
		}
	}

	assertSourceOmits(t, filepath.Join(root, "protocol", "v2", "jsonrpc.go"),
		"agent.file", "FileOperation", "FileResult")
	assertSourceOmits(t, filepath.Join(root, "server", "websocket.go"),
		"MethodAgentFile", "handleFileOperation", `"file"`)
}

func assertSourceOmits(t *testing.T, path string, forbidden ...string) {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	for _, token := range forbidden {
		if strings.Contains(string(content), token) {
			t.Errorf("%s still contains remote-control token %q", path, token)
		}
	}
}
