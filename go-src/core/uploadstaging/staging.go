// Package uploadstaging picks disk paths for large CMS/deploy uploads.
// Server /tmp is a small tmpfs; pub-tool scripts use /home/ec2-user/staging — runtime uploads follow the same rule.
package uploadstaging

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	envStagingRoot      = "SARA_UPLOAD_STAGING_ROOT"
	defaultLinuxStaging = "/home/ec2-user/staging"

	// SubOfficialSiteFileUpload holds chunked official-site / third-pay file upload sessions.
	SubOfficialSiteFileUpload = "official-site-file-upload"
	// SubCMSDeployZip holds incoming deploy zip files from CMS multipart uploads.
	SubCMSDeployZip = "cms-deploy-zip"
	// SubDBBackup holds mysqldump gzip files before upload to object storage.
	SubDBBackup = "db-backup"
)

// Root returns the disk staging root (same default as pub-tool REMOTE_STAGE).
func Root() (string, error) {
	var root string
	if v := strings.TrimSpace(os.Getenv(envStagingRoot)); v != "" {
		root = filepath.Clean(v)
	} else if runtime.GOOS == "linux" {
		root = defaultLinuxStaging
	} else {
		cache, err := os.UserCacheDir()
		if err != nil || strings.TrimSpace(cache) == "" {
			cache = "."
		}
		root = filepath.Join(cache, "sara-staging")
	}
	if err := os.MkdirAll(root, 0755); err != nil {
		return "", err
	}
	return root, nil
}

// Subdir ensures root/subdir exists and returns its path.
func Subdir(name string) (string, error) {
	root, err := Root()
	if err != nil {
		return "", err
	}
	name = strings.TrimSpace(name)
	if name == "" || strings.Contains(name, "..") {
		return "", os.ErrInvalid
	}
	p := filepath.Join(root, name)
	if err := os.MkdirAll(p, 0700); err != nil {
		return "", err
	}
	return p, nil
}

// CreateTempFile creates a file under root/subdir (never under tmpfs /tmp on Linux servers).
func CreateTempFile(subdir, pattern string) (*os.File, error) {
	dir, err := Subdir(subdir)
	if err != nil {
		return nil, err
	}
	return os.CreateTemp(dir, pattern)
}
