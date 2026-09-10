package recharge

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func haiPayNormalizeKeyMaterial(raw string) string {
	s := strings.TrimSpace(raw)
	s = strings.ReplaceAll(s, "\r", "")
	return s
}

func haiPayParsePrivateKey(raw string) (*rsa.PrivateKey, error) {
	raw = haiPayNormalizeKeyMaterial(raw)
	if raw == "" {
		return nil, fmt.Errorf("empty private key")
	}
	if block, _ := pem.Decode([]byte(raw)); block != nil {
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			if pkcs1, err2 := x509.ParsePKCS1PrivateKey(block.Bytes); err2 == nil {
				return pkcs1, nil
			}
			return nil, err
		}
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("not rsa private key")
		}
		return rsaKey, nil
	}
	der, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(raw, "\n", ""))
	if err != nil {
		return nil, err
	}
	key, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		if pkcs1, err2 := x509.ParsePKCS1PrivateKey(der); err2 == nil {
			return pkcs1, nil
		}
		return nil, err
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not rsa private key")
	}
	return rsaKey, nil
}

func haiPayParsePublicKey(raw string) (*rsa.PublicKey, error) {
	raw = haiPayNormalizeKeyMaterial(raw)
	if raw == "" {
		return nil, fmt.Errorf("empty public key")
	}
	if block, _ := pem.Decode([]byte(raw)); block != nil {
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaKey, ok := key.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("not rsa public key")
		}
		return rsaKey, nil
	}
	der, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(raw, "\n", ""))
	if err != nil {
		return nil, err
	}
	key, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not rsa public key")
	}
	return rsaKey, nil
}

func haiPayValueToSignString(v any) (string, bool) {
	if v == nil {
		return "", false
	}
	switch t := v.(type) {
	case string:
		if t == "" {
			return "", false
		}
		return t, true
	case float64:
		// JSON numbers decode as float64; keep integer-looking values without .0
		if t == float64(int64(t)) {
			return strconv.FormatInt(int64(t), 10), true
		}
		return strconv.FormatFloat(t, 'f', -1, 64), true
	case float32:
		f := float64(t)
		if f == float64(int64(f)) {
			return strconv.FormatInt(int64(f), 10), true
		}
		return strconv.FormatFloat(f, 'f', -1, 64), true
	case int:
		return strconv.Itoa(t), true
	case int64:
		return strconv.FormatInt(t, 10), true
	case int32:
		return strconv.FormatInt(int64(t), 10), true
	case uint64:
		return strconv.FormatUint(t, 10), true
	case bool:
		if t {
			return "true", true
		}
		return "false", true
	default:
		s := strings.TrimSpace(fmt.Sprint(t))
		if s == "" || s == "<nil>" {
			return "", false
		}
		return s, true
	}
}

func haiPayBuildSignContent(params map[string]any, merchantSecretKey string) string {
	keys := make([]string, 0, len(params))
	for k, v := range params {
		if k == "sign" || k == "sign_type" {
			continue
		}
		if _, ok := haiPayValueToSignString(v); !ok {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys)+1)
	for _, k := range keys {
		val, _ := haiPayValueToSignString(params[k])
		parts = append(parts, k+"="+val)
	}
	parts = append(parts, "key="+merchantSecretKey)
	return strings.Join(parts, "&")
}

func haiPaySign(params map[string]any, merchantSecretKey, privateKeyPEM string) (string, error) {
	content := haiPayBuildSignContent(params, merchantSecretKey)
	priv, err := haiPayParsePrivateKey(privateKeyPEM)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(content))
	sig, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, sum[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

func haiPayVerify(params map[string]any, merchantSecretKey, publicKeyPEM, signBase64 string) error {
	signBase64 = strings.TrimSpace(signBase64)
	if signBase64 == "" {
		return fmt.Errorf("empty sign")
	}
	content := haiPayBuildSignContent(params, merchantSecretKey)
	pub, err := haiPayParsePublicKey(publicKeyPEM)
	if err != nil {
		return err
	}
	sig, err := base64.StdEncoding.DecodeString(signBase64)
	if err != nil {
		return err
	}
	sum := sha256.Sum256([]byte(content))
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, sum[:], sig)
}
