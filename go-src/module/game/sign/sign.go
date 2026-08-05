package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

const (
	FieldSign           = "sign"
	HeaderOperatorToken = "operator_token"
	HeaderTimestamp     = "timestamp"

	HTTPHeaderOperatorToken = "X-Operator-Token"
	HTTPHeaderTimestamp     = "X-Timestamp"
	HTTPHeaderSign          = "X-Sign"
)

// Build 根据参与签名的参数生成 X-Sign(小写十六进制).
func Build(params map[string]string, secretKey string) string {
	secretKey = strings.TrimSpace(secretKey)
	if secretKey == "" || len(params) == 0 {
		return ""
	}
	signStr := BuildString(params, secretKey)
	mac := hmac.New(sha256.New, []byte(secretKey))
	_, _ = mac.Write([]byte(signStr))
	return hex.EncodeToString(mac.Sum(nil))
}

// Verify 校验签名是否与参数匹配.
func Verify(params map[string]string, sign, secretKey string) bool {
	sign = strings.TrimSpace(sign)
	if sign == "" {
		return false
	}
	expected := Build(params, secretKey)
	return hmac.Equal([]byte(expected), []byte(strings.ToLower(sign)))
}

// BuildString 返回 HMAC 前的待签名字符串,便于排查.
func BuildString(params map[string]string, secretKey string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "" || k == FieldSign {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+params[k])
	}
	return strings.Join(parts, "&") + "&secret=" + secretKey
}

// MergeParams 合并 query、body 与 Header 中参与签名的字段.
func MergeParams(query map[string]string, body map[string]string, operatorToken, timestamp string) map[string]string {
	params := make(map[string]string)
	for k, v := range query {
		if k == "" || k == FieldSign {
			continue
		}
		params[k] = v
	}
	for k, v := range body {
		if k == "" || k == FieldSign {
			continue
		}
		params[k] = v
	}
	if operatorToken = strings.TrimSpace(operatorToken); operatorToken != "" {
		params[HeaderOperatorToken] = operatorToken
	}
	if timestamp = strings.TrimSpace(timestamp); timestamp != "" {
		params[HeaderTimestamp] = timestamp
	}
	return params
}

// ParseQuery 解析 URL query 参数.
func ParseQuery(values url.Values) map[string]string {
	params := make(map[string]string)
	for k, items := range values {
		if len(items) == 0 {
			continue
		}
		params[k] = items[len(items)-1]
	}
	return params
}

// ParseJSONBody 解析 JSON body 顶层字段为 key-value.
func ParseJSONBody(raw []byte) map[string]string {
	if len(raw) == 0 {
		return map[string]string{}
	}
	var body map[string]any
	if err := json.Unmarshal(raw, &body); err != nil {
		return map[string]string{}
	}
	return flattenJSONBody(body)
}

func flattenJSONBody(body map[string]any) map[string]string {
	params := make(map[string]string, len(body))
	for k, v := range body {
		if k == "" || k == FieldSign {
			continue
		}
		params[k] = flattenValue(v)
	}
	return params
}

func flattenValue(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case json.Number:
		return x.String()
	case float64:
		if x == float64(int64(x)) {
			return strconv.FormatInt(int64(x), 10)
		}
		return strconv.FormatFloat(x, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(x)
	default:
		b, err := json.Marshal(x)
		if err != nil {
			return fmt.Sprint(x)
		}
		return string(b)
	}
}
