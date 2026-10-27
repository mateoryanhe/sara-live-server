package auditcode

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/util/guid"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/auditcodedto"
	"xr-game-server/entity/sys"
	"xr-game-server/errercode"
)

const (
	// legacyDefaultAuditCode 库中尚无记录时校验兼容旧硬编码主播码.
	legacyDefaultAuditCode = "971200"
	defaultAuditCodeLen    = 8
	minAuditCodeLen        = 4
	maxAuditCodeLen        = 32
)

func GetAuditCodeCfg(_ context.Context, _ *auditcodedto.GetAuditCodeCfgReq) (*auditcodedto.GetAuditCodeCfgRes, error) {
	return &auditcodedto.GetAuditCodeCfgRes{Cfg: toCfgItem(cfgdao.LoadAuditCodeCfg())}, nil
}

func GenerateAuditCodeCfg(_ context.Context, req *auditcodedto.GenerateAuditCodeCfgReq) (*auditcodedto.GenerateAuditCodeCfgRes, error) {
	length := defaultAuditCodeLen
	if req != nil && req.Length > 0 {
		length = req.Length
	}
	if length < minAuditCodeLen || length > maxAuditCodeLen {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	code := newAuditCode(length)
	existing := cfgdao.LoadAuditCodeCfg()
	row := &entity.AuditCodeCfg{Code: code}
	now := time.Now()
	if existing != nil {
		row.ID = existing.ID
		row.CreatedAt = existing.CreatedAt
	} else {
		row.CreatedAt = now
	}
	row.UpdatedAt = now
	if err := cfgdao.SaveAuditCodeCfg(row); err != nil {
		return nil, err
	}
	return &auditcodedto.GenerateAuditCodeCfgRes{
		Success: true,
		Cfg:     toCfgItem(row),
	}, nil
}

// newAuditCode 生成小写、无下划线的 guid 片段,长度为 length.
func newAuditCode(length int) string {
	var b strings.Builder
	b.Grow(length)
	for b.Len() < length {
		part := strings.ToLower(strings.ReplaceAll(guid.S(), "_", ""))
		b.WriteString(part)
	}
	return b.String()[:length]
}

// CurrentAuditCode 每次从库读取当前审核码(无缓存);库空则兼容旧默认码.
func CurrentAuditCode() string {
	row := cfgdao.LoadAuditCodeCfg()
	if row == nil {
		return legacyDefaultAuditCode
	}
	code := strings.TrimSpace(row.Code)
	if code == "" {
		return legacyDefaultAuditCode
	}
	return code
}

func toCfgItem(row *entity.AuditCodeCfg) *auditcodedto.AuditCodeCfgItem {
	if row == nil {
		return nil
	}
	return &auditcodedto.AuditCodeCfgItem{
		ID:        strconv.FormatUint(row.ID, 10),
		Code:      row.Code,
		CreatedAt: formatTime(row.CreatedAt),
		UpdatedAt: formatTime(row.UpdatedAt),
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
