package auditcode

import (
	"context"
	"crypto/rand"
	"strconv"
	"strings"
	"time"

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
	code, err := newAuditCode(length)
	if err != nil {
		return nil, err
	}
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

// auditCodeAlphabet 数字+字母(去掉易混淆 0/O、1/l/I 等),便于 App 手工输入.
const auditCodeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789"

// newAuditCode 密码学随机,长度为 length.
func newAuditCode(length int) (string, error) {
	if length <= 0 {
		return "", errercode.CreateCode(errercode.InvalidParam)
	}
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	for i := range buf {
		buf[i] = auditCodeAlphabet[int(buf[i])%len(auditCodeAlphabet)]
	}
	return string(buf), nil
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
