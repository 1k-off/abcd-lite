package deployment

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Helper to create a temp dir with files
func createTempDirWithFiles(t *testing.T, files []string) string {
	dir := t.TempDir()
	for _, f := range files {
		fpath := filepath.Join(dir, f)
		os.WriteFile(fpath, []byte("test"), 0644)
	}
	return dir
}

func TestClean_RemovesUnexcludedFiles(t *testing.T) {
	dir := createTempDirWithFiles(t, []string{"a.txt", "b.txt", "c.txt"})
	exclude := []string{"b.txt"}

	err := clean(dir, exclude)
	if err != nil {
		t.Fatalf("clean failed: %v", err)
	}

	files, _ := os.ReadDir(dir)
	var names []string
	for _, f := range files {
		names = append(names, f.Name())
	}
	if len(names) != 1 || names[0] != "b.txt" {
		t.Errorf("expected only b.txt to remain, got %v", names)
	}
}

func TestClean_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	err := clean(dir, []string{})
	if err != nil {
		t.Fatalf("clean failed on empty dir: %v", err)
	}
}

func TestClean_ReadDirError(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "nonexistent-dir-for-test")
	err := clean(dir, []string{})
	if err == nil || !strings.Contains(err.Error(), "failed to read destination") {
		t.Errorf("expected read error, got %v", err)
	}
}

func TestDeploy_InvalidOptions(t *testing.T) {
	info := PackageInfo{}
	err := Deploy(Options{}, info)
	if err == nil || !strings.Contains(err.Error(), "destination is required") {
		t.Errorf("expected destination error, got %v", err)
	}

	err = Deploy(Options{Destination: "here"}, info)
	if err == nil || !strings.Contains(err.Error(), "oras pull failed") {
		t.Errorf("expected oras pull failure, got %v", err)
	}
}

func TestClean_ExcludeWithTrailingSlash(t *testing.T) {
	dir := createTempDirWithFiles(t, []string{"dir", "file.txt"})
	exclude := []string{"dir/"}

	err := clean(dir, exclude)
	if err != nil {
		t.Fatalf("clean failed: %v", err)
	}

	files, _ := os.ReadDir(dir)
	var names []string
	for _, f := range files {
		names = append(names, f.Name())
	}
	if len(names) != 1 || names[0] != "dir" {
		t.Errorf("expected only dir to remain, got %v", names)
	}
}

func TestClean_ExcludeWithoutTrailingSlash(t *testing.T) {
	dir := createTempDirWithFiles(t, []string{"dir", "file.txt"})
	exclude := []string{"dir"}

	err := clean(dir, exclude)
	if err != nil {
		t.Fatalf("clean failed: %v", err)
	}

	files, _ := os.ReadDir(dir)
	var names []string
	for _, f := range files {
		names = append(names, f.Name())
	}
	if len(names) != 1 || names[0] != "dir" {
		t.Errorf("expected only dir to remain, got %v", names)
	}
}
