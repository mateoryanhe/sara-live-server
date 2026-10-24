package sysdto

import (
	"encoding/json"
	"testing"
)

func TestSysCfgRespThirdPayFieldNames(t *testing.T) {
	data, err := json.Marshal(SysCfgResp{
		T:               "third-pay.bigtktool.shop",
		TVisable:        true,
		W:               "withdraw.example.com",
		WVisable:        true,
		WAndroidVisable: true,
	})
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if got := payload["t"]; got != "third-pay.bigtktool.shop" {
		t.Fatalf("payload[t] = %v", got)
	}
	if got := payload["tVisable"]; got != true {
		t.Fatalf("payload[tVisable] = %v", got)
	}
	if got := payload["w"]; got != "withdraw.example.com" {
		t.Fatalf("payload[w] = %v", got)
	}
	if got := payload["wVisable"]; got != true {
		t.Fatalf("payload[wVisable] = %v", got)
	}
	if got := payload["wAndroidVisable"]; got != true {
		t.Fatalf("payload[wAndroidVisable] = %v", got)
	}
}
