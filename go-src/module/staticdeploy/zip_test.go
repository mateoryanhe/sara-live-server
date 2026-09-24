package staticdeploy

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestDeployZipReplacesExistingSite(t *testing.T) {
	root := t.TempDir()
	destRoot := filepath.Join(root, "site")
	if err := os.MkdirAll(destRoot, 0755); err != nil {
		t.Fatal(err)
	}
	oldPath := filepath.Join(destRoot, "old-only.js")
	if err := os.WriteFile(oldPath, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}

	zipPath := filepath.Join(root, "release.zip")
	writeTestZip(t, zipPath, map[string]string{
		"index.html":    "new index",
		"assets/app.js": "new app",
	})

	fileCount, _, err := DeployZip(zipPath, destRoot)
	if err != nil {
		t.Fatalf("DeployZip() error = %v", err)
	}
	if fileCount != 2 {
		t.Fatalf("fileCount = %d, want 2", fileCount)
	}
	if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
		t.Fatalf("obsolete file still exists, stat error = %v", err)
	}
	content, err := os.ReadFile(filepath.Join(destRoot, "assets", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "new app" {
		t.Fatalf("deployed content = %q, want %q", string(content), "new app")
	}
}

func TestDeployZipKeepsExistingSiteWhenPackageIsInvalid(t *testing.T) {
	root := t.TempDir()
	destRoot := filepath.Join(root, "site")
	if err := os.MkdirAll(destRoot, 0755); err != nil {
		t.Fatal(err)
	}
	oldPath := filepath.Join(destRoot, "index.html")
	if err := os.WriteFile(oldPath, []byte("current site"), 0644); err != nil {
		t.Fatal(err)
	}

	zipPath := filepath.Join(root, "broken.zip")
	if err := os.WriteFile(zipPath, []byte("not a zip"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := DeployZip(zipPath, destRoot); err == nil {
		t.Fatal("DeployZip() error = nil, want invalid zip error")
	}
	content, err := os.ReadFile(oldPath)
	if err != nil {
		t.Fatalf("current site was removed: %v", err)
	}
	if string(content) != "current site" {
		t.Fatalf("current site content = %q, want unchanged content", string(content))
	}
}

func writeTestZip(t *testing.T, zipPath string, files map[string]string) {
	t.Helper()
	output, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	writer := zip.NewWriter(output)
	for name, content := range files {
		entry, err := writer.Create(name)
		if err != nil {
			_ = output.Close()
			t.Fatal(err)
		}
		if _, err := entry.Write([]byte(content)); err != nil {
			_ = output.Close()
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		_ = output.Close()
		t.Fatal(err)
	}
	if err := output.Close(); err != nil {
		t.Fatal(err)
	}
}
