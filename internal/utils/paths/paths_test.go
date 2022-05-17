package paths_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/reaper47/recipya/internal/utils/paths"
)

func TestData(t *testing.T) {
	got := paths.Data()
	want := filepath.Join(exe(), "data")
	assertPaths(t, got, want)
}

func TestImages(t *testing.T) {
	got := paths.Images()
	want := filepath.Join(exe(), "data", "img")
	assertPaths(t, got, want)
}

func exe() string {
	exe, _ := os.Executable()
	return filepath.Dir(exe)
}

func assertPaths(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q but want %q", got, want)
	}
}
