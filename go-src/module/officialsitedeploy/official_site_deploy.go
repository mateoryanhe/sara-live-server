package officialsitedeploy

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/cfg"
	"xr-game-server/dto/officialsitedeploydto"
	"xr-game-server/errercode"
	"xr-game-server/module/domainsite"
	"xr-game-server/module/staticdeploy"
)

type deployTarget struct {
	siteKey     string
	prefix      string
	tmpPattern  string
	acceptExt   string
	apkSubDir   string
	defaultPath string
}

var (
	officialSiteTarget = deployTarget{
		siteKey:     officialsitedeploydto.OfficialSiteSiteKey,
		prefix:      officialsitedeploydto.OfficialSiteStaticPrefix,
		tmpPattern:  "official-site-deploy-*.zip",
		acceptExt:   ".zip,.apk",
		apkSubDir:   "assets",
		defaultPath: officialsitedeploydto.DefaultOfficialSiteDeployPath,
	}
	thirdPayOfficialSiteTarget = deployTarget{
		siteKey:     officialsitedeploydto.ThirdPayOfficialSiteSiteKey,
		prefix:      officialsitedeploydto.ThirdPayOfficialSiteStaticPrefix,
		tmpPattern:  "third-pay-official-site-deploy-*.zip",
		acceptExt:   ".zip,.apk",
		apkSubDir:   "assets",
		defaultPath: officialsitedeploydto.DefaultThirdPayOfficialSitePath,
	}
)

func GetOfficialSiteDeployInfo(ctx context.Context, _ *officialsitedeploydto.GetOfficialSiteDeployInfoReq) (*officialsitedeploydto.GetOfficialSiteDeployInfoRes, error) {
	info, err := getDeployInfo(ctx, officialSiteTarget)
	if err != nil {
		return nil, err
	}
	return &officialsitedeploydto.GetOfficialSiteDeployInfoRes{Info: info}, nil
}

func SaveOfficialSiteDeployCfg(ctx context.Context, req *officialsitedeploydto.SaveOfficialSiteDeployCfgReq) (*officialsitedeploydto.SaveOfficialSiteDeployCfgRes, error) {
	domain := strings.TrimSpace(req.Domain)
	deployPath := strings.TrimSpace(req.DeployPath)
	if domain == "" || deployPath == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if _, err := domainsite.SaveStaticSiteMapping(
		ctx,
		officialsitedeploydto.OfficialSiteSiteKey,
		officialsitedeploydto.OfficialSiteStaticPrefix,
		domain,
		deployPath,
	); err != nil {
		return nil, err
	}
	return &officialsitedeploydto.SaveOfficialSiteDeployCfgRes{Success: true}, nil
}

func GetThirdPayOfficialSiteDeployInfo(ctx context.Context, _ *officialsitedeploydto.GetThirdPayOfficialSiteDeployInfoReq) (*officialsitedeploydto.GetThirdPayOfficialSiteDeployInfoRes, error) {
	info, err := getDeployInfo(ctx, thirdPayOfficialSiteTarget)
	if err != nil {
		return nil, err
	}
	return &officialsitedeploydto.GetThirdPayOfficialSiteDeployInfoRes{Info: info}, nil
}

func SaveThirdPayOfficialSiteDeployCfg(ctx context.Context, req *officialsitedeploydto.SaveThirdPayOfficialSiteDeployCfgReq) (*officialsitedeploydto.SaveThirdPayOfficialSiteDeployCfgRes, error) {
	domain := strings.TrimSpace(req.Domain)
	deployPath := strings.TrimSpace(req.DeployPath)
	if domain == "" || deployPath == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if _, err := domainsite.SaveStaticSiteMapping(
		ctx,
		officialsitedeploydto.ThirdPayOfficialSiteSiteKey,
		officialsitedeploydto.ThirdPayOfficialSiteStaticPrefix,
		domain,
		deployPath,
	); err != nil {
		return nil, err
	}
	return &officialsitedeploydto.SaveThirdPayOfficialSiteDeployCfgRes{Success: true}, nil
}

func getDeployInfo(ctx context.Context, target deployTarget) (*officialsitedeploydto.SiteDeployInfoItem, error) {
	deployPath, err := getDeployDir(target)
	if err != nil {
		return nil, err
	}
	return &officialsitedeploydto.SiteDeployInfoItem{
		Domain:       cfg.GetStaticSiteDomains(target.prefix),
		UrlPrefix:    target.prefix,
		DeployPath:   deployPath,
		AcceptExt:    target.acceptExt,
		LastUploadAt: domainsite.GetLastUploadAt(ctx, target.siteKey),
	}, nil
}

func DeployOfficialSiteZipFromRequest(r *ghttp.Request) (*officialsitedeploydto.DeploySiteZipRes, error) {
	return deployZipFromRequest(r, officialSiteTarget)
}

func DeployThirdPayOfficialSiteZipFromRequest(r *ghttp.Request) (*officialsitedeploydto.DeploySiteZipRes, error) {
	return deployZipFromRequest(r, thirdPayOfficialSiteTarget)
}

func deployZipFromRequest(r *ghttp.Request, target deployTarget) (*officialsitedeploydto.DeploySiteZipRes, error) {
	if r == nil || r.Request == nil {
		return nil, errors.New("upload file is empty")
	}
	deployDir, err := getDeployDir(target)
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
		tmpFile, err := os.CreateTemp("", target.tmpPattern)
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
	if _, err = domainsite.RecordDeploySuccess(r.Context(), target.siteKey); err != nil {
		return nil, fmt.Errorf("files deployed but record upload time failed: %w", err)
	}
	return &officialsitedeploydto.DeploySiteZipRes{
		FileCount:  fileCount,
		DirCount:   dirCount,
		DeployPath: deployDir,
		UrlPrefix:  target.prefix,
	}, nil
}

func getDeployDir(target deployTarget) (string, error) {
	root := strings.TrimSpace(cfg.GetStaticPathRoot(target.prefix))
	if root == "" {
		root = target.defaultPath
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
