//go:build manual

package recharge

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"xr-game-server/constants/country"
)

type haipaySmokeCfg struct {
	ApiHost            string `json:"apiHost"`
	MerchantSecretKey  string `json:"merchantSecretKey"`
	MerchantPrivateKey string `json:"merchantPrivateKey"`
	CallbackBaseUrl    string `json:"callbackBaseUrl"`
	ReturnUrl          string `json:"returnUrl"`
	FailReturnUrl      string `json:"failReturnUrl"`
	UseProdAppID       int    `json:"useProdAppID"`
}

// 本地探活正式 HaiPay：从正式服导出 haipay_cfgs 为 JSON，勿提交该文件。
//
//	go test ./module/recharge -tags manual -run TestManualHaiPayGlobalCashierApply -count=1 -v \
//	  -cfg C:/Users/hw/AppData/Local/Temp/haipay_prod_cfg.json
var smokeCfgPath = flag.String("cfg", "", "path to haipay_cfgs JSON (required)")
var smokeForce = flag.Bool("force", false, "skip local private-key parse; still POST HaiPay (bogus sign if local sign fails)")

func TestManualHaiPayGlobalCashierApply(t *testing.T) {
	runManualHaiPayGlobalCashierApply(t, false)
}

func TestManualHaiPayGlobalCashierApplyForce(t *testing.T) {
	runManualHaiPayGlobalCashierApply(t, true)
}

func runManualHaiPayGlobalCashierApply(t *testing.T, force bool) {
	t.Helper()
	flag.Parse()
	if force {
		*smokeForce = true
	}
	if *smokeCfgPath == "" {
		t.Fatal("pass -cfg with prod haipay_cfgs JSON")
	}
	raw, err := os.ReadFile(*smokeCfgPath)
	if err != nil {
		t.Fatal(err)
	}
	var cfg haipaySmokeCfg
	if err = json.Unmarshal(raw, &cfg); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(cfg.MerchantPrivateKey) == "" || strings.TrimSpace(cfg.MerchantSecretKey) == "" {
		t.Fatal("incomplete cfg")
	}

	ctx := context.Background()
	orderID := fmt.Sprintf("smoke%d", time.Now().UnixNano())
	userID := uint64(1000039)
	amount := "17717.04"
	region := "ID"
	currency := "IDR"

	appID := country.HaiPayCashierAppID(cfg.UseProdAppID != 0)
	returnURL := strings.TrimSpace(cfg.ReturnUrl)
	if returnURL == "" {
		returnURL = cfg.CallbackBaseUrl
	}
	failURL := strings.TrimSpace(cfg.FailReturnUrl)
	if failURL == "" {
		failURL = returnURL
	}
	body := map[string]any{
		"appId":           appID,
		"orderId":         orderID,
		"name":            "Smoke User",
		"email":           "smoke@noreply.local",
		"amount":          amount,
		"currency":        currency,
		"callBackUrl":     returnURL,
		"callBackFailUrl": failURL,
		"notifyUrl":       strings.TrimRight(cfg.CallbackBaseUrl, "/") + haiPayCollectNotifyPath,
		"subject":         "Recharge",
		"region":          region,
		"partnerUserId":   strconv.FormatUint(userID, 10),
		"body":            "order:" + orderID,
	}

	if !*smokeForce {
		if _, err = haiPayParsePrivateKey(cfg.MerchantPrivateKey); err != nil {
			t.Fatalf("private key parse: %v", err)
		}
	} else {
		if _, parseErr := haiPayParsePrivateKey(cfg.MerchantPrivateKey); parseErr != nil {
			t.Logf("force: local private key parse failed (ignored): %v", parseErr)
		}
	}
	sign, err := haiPaySign(ctx, body, cfg.MerchantSecretKey, cfg.MerchantPrivateKey)
	if err != nil {
		if !*smokeForce {
			t.Fatalf("sign: %v", err)
		}
		t.Logf("force: local sign failed, posting with placeholder sign: %v", err)
		body["sign"] = "FORCE_INVALID_SIGN_SMOKE"
	} else {
		body["sign"] = sign
	}

	var res haiPayApplyAPIRes
	url := strings.TrimRight(cfg.ApiHost, "/") + haiPayGlobalCollectApplyPath
	t.Logf("POST %s orderId=%s appId=%d force=%v", url, orderID, appID, *smokeForce)
	if err = haiPayPostJSON(ctx, url, body, &res); err != nil {
		t.Fatal(err)
	}
	t.Logf("response status=%s error=%s msg=%s", res.Status, res.Error, res.Msg)
	if res.Data != nil {
		t.Logf("data orderNo=%s payUrl=%s", res.Data.OrderNo, res.Data.PayUrl)
	}
	if *smokeForce {
		return
	}
	if res.Status != "1" {
		t.Fatalf("apply failed")
	}
	if res.Data == nil || strings.TrimSpace(res.Data.PayUrl) == "" {
		t.Fatalf("empty payUrl data=%+v", res.Data)
	}
	t.Logf("payUrl=%s orderNo=%s", res.Data.PayUrl, res.Data.OrderNo)
}
