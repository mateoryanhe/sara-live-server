package uploadstaging

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSubdirUsesEnvRoot(t *testing.T) {
	root := t.TempDir()
	t.Setenv(envStagingRoot, root)

	got, err := Subdir("test-sub")
	if err != nil {
		t.Fatalf("Subdir: %v", err)
	}
	want := filepath.Join(root, "test-sub")
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	if _, err := os.Stat(want); err != nil {
		t.Fatalf("stat subdir: %v", err)
	}
}

func TestCreateTempFileUnderSubdir(t *testing.T) {
	root := t.TempDir()
	t.Setenv(envStagingRoot, root)

	f, err := CreateTempFile(SubCMSDeployZip, "unit-*.zip")
	if err != nil {
		t.Fatalf("CreateTempFile: %v", err)
	}
	path := f.Name()
	_ = f.Close()
	_ = os.Remove(path)

	if !strings.HasPrefix(path, filepath.Join(root, SubCMSDeployZip)+string(os.PathSeparator)) {
		t.Fatalf("temp not under staging subdir: %s", path)
	}
}
