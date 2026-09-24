package coinmerchantdeploy

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/cfg"
	"xr-game-server/dto/coinmerchantdeploydto"
	"xr-game-server/module/domainsite"
	"xr-game-server/module/staticdeploy"
)

func DeployZipFromRequest(r *ghttp.Request) (*coinmerchantdeploydto.DeployCoinMerchantZipRes, error) {
	if r == nil || r.Request == nil {
		return nil, errors.New("upload file is empty")
	}
	deployDir, err := getDeployDir()
	if err != nil {
		return nil, err
	}

	reader, err := r.Request.MultipartReader()
	if err != nil {
		return nil, mapUploadReadErr(err)
	}

	var zipPath string
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, mapUploadReadErr(err)
		}
		if part.FormName() != "file" {
			part.Close()
			continue
		}
		ext := strings.ToLower(filepath.Ext(part.FileName()))
		if ext != ".zip" {
			part.Close()
			return nil, fmt.Errorf("file ext not allowed: %s", ext)
		}
		tmpFile, err := os.CreateTemp("", "coin-merchant-deploy-*.zip")
		if err != nil {
			part.Close()
			return nil, err
		}
		zipPath = tmpFile.Name()
		_, copyErr := io.Copy(tmpFile, part)
		part.Close()
		closeErr := tmpFile.Close()
		if copyErr != nil {
			os.Remove(zipPath)
			return nil, mapUploadReadErr(copyErr)
		}
		if closeErr != nil {
			os.Remove(zipPath)
			return nil, closeErr
		}
		break
	}
	if zipPath == "" {
		return nil, errors.New("upload file is empty")
	}
	defer os.Remove(zipPath)

	fileCount, dirCount, err := staticdeploy.DeployZip(zipPath, deployDir)
	if err != nil {
		return nil, err
	}
	if _, err = domainsite.RecordDeploySuccess(r.Context(), coinmerchantdeploydto.CoinMerchantSiteKey); err != nil {
		return nil, fmt.Errorf("files deployed but record upload time failed: %w", err)
	}
	return &coinmerchantdeploydto.DeployCoinMerchantZipRes{
		FileCount:  fileCount,
		DirCount:   dirCount,
		DeployPath: deployDir,
		UrlPrefix:  coinmerchantdeploydto.CoinMerchantStaticPrefix,
	}, nil
}

func getDeployDir() (string, error) {
	root := strings.TrimSpace(cfg.GetStaticPathRoot(coinmerchantdeploydto.CoinMerchantStaticPrefix))
	if root == "" {
		root = coinmerchantdeploydto.DefaultCoinMerchantDeployPath
	}
	if err := os.MkdirAll(root, 0755); err != nil {
		return "", fmt.Errorf("create deploy dir %s: %w", root, err)
	}
	return root, nil
}

func mapUploadReadErr(err error) error {
	if err == nil {
		return nil
	}
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		return errors.New("upload request incomplete")
	}
	if strings.Contains(strings.ToLower(err.Error()), "unexpected eof") {
		return errors.New("upload request incomplete")
	}
	return err
}
