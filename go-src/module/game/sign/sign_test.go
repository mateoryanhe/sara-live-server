package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
)

func TestBuildMatchesNodeExample(t *testing.T) {
	secretKey := "test-secret"
	params := map[string]string{
		"foo": "bar",
		"baz": "123",
	}
	signStr := BuildString(params, secretKey)
	mac := hmac.New(sha256.New, []byte(secretKey))
	_, _ = mac.Write([]byte(signStr))
	expected := hex.EncodeToString(mac.Sum(nil))

	got := Build(params, secretKey)
	if got != expected {
		t.Fatalf("Build() = %q, want %q", got, expected)
	}
	if !Verify(params, expected, secretKey) {
		t.Fatal("Verify() should accept valid sign")
	}
}

func TestVerifyRejectsInvalidSign(t *testing.T) {
	params := map[string]string{"foo": "bar"}
	if Verify(params, "invalid", "test-secret") {
		t.Fatal("Verify() should reject invalid sign")
	}
}

func TestMergeParamsIncludesHeadersAndExcludesSign(t *testing.T) {
	params := MergeParams(
		map[string]string{"a": "1"},
		map[string]string{"b": "2", FieldSign: "should-ignore"},
		"token-abc",
		"1700000000",
	)
	if params["a"] != "1" || params["b"] != "2" {
		t.Fatalf("unexpected query/body params: %#v", params)
	}
	if params[HeaderOperatorToken] != "token-abc" {
		t.Fatalf("operator_token = %q", params[HeaderOperatorToken])
	}
	if params[HeaderTimestamp] != "1700000000" {
		t.Fatalf("timestamp = %q", params[HeaderTimestamp])
	}
	if _, ok := params[FieldSign]; ok {
		t.Fatal("sign should not be included in params")
	}
}

func TestBuildStringSorted(t *testing.T) {
	got := BuildString(map[string]string{
		"z": "last",
		"a": "first",
	}, "secret")
	if !strings.HasPrefix(got, "a=first&z=last&secret=secret") {
		t.Fatalf("unexpected sign string: %q", got)
	}
}
