package entity

import (
	"reflect"
	"strings"
	"testing"
)

func TestAnchorSalarySocialSharePercentColumnMapping(t *testing.T) {
	field, ok := reflect.TypeOf(AnchorSalarySocialShareCfg{}).FieldByName("AnchorSocialSharePercent")
	if !ok {
		t.Fatal("AnchorSocialSharePercent field is missing")
	}
	if got := strings.Split(field.Tag.Get("orm"), ",")[0]; got != "social_share_percent" {
		t.Fatalf("GoFrame ORM column = %q, want social_share_percent", got)
	}
	if got := strings.Split(field.Tag.Get("gorm"), ";")[0]; got != "column:social_share_percent" {
		t.Fatalf("GORM column = %q, want column:social_share_percent", got)
	}
}
