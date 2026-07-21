package phoneutil

import "strings"

func NormalizeAreaCode(areaCode string) string {
	areaCode = strings.TrimSpace(areaCode)
	return strings.TrimPrefix(areaCode, "+")
}

func UniqueKey(areaCode, phone string) string {
	return NormalizeAreaCode(areaCode) + ":" + strings.TrimSpace(phone)
}
