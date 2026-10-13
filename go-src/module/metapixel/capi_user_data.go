package metapixel

import (
	"regexp"
	"strings"
	"unicode"

	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dao/userloginlocationdao"
	"xr-game-server/gameevent"
)

const metaEmailAccountChannel uint = 7 // 与 auth.EmailChannel 一致

var metaNonDigit = regexp.MustCompile(`[^0-9]`)

// capiProfileHint 单次事件携带的明文匹配字段(优先于库内资料).
type capiProfileHint struct {
	Email string
	Phone string
	Name  string
}

func buildCAPIUserData(userId uint64, attr gameevent.RegisterAttribution, hint capiProfileHint) capiUserData {
	out := capiUserData{}
	if extID := hashMetaExternalID(userId); extID != "" {
		out.ExternalID = []string{extID}
	}
	applyProfileHashes(userId, hint, &out)
	if ip := strings.TrimSpace(attr.ClientIPAddress); ip != "" {
		out.ClientIPAddress = ip
	}
	if ua := strings.TrimSpace(attr.ClientUserAgent); ua != "" {
		out.ClientUserAgent = ua
	}
	if fbp := strings.TrimSpace(attr.MetaFbp); fbp != "" {
		out.Fbp = fbp
	}
	if fbc := strings.TrimSpace(attr.MetaFbc); fbc != "" {
		out.Fbc = fbc
	}
	return out
}

func buildCAPIUserDataForUser(userId uint64) capiUserData {
	return buildCAPIUserData(userId, gameevent.RegisterAttribution{}, capiProfileHint{})
}

func checkoutHintFromEvent(c gameevent.RechargeCheckoutHint) capiProfileHint {
	return capiProfileHint{
		Email: strings.TrimSpace(c.Email),
		Phone: strings.TrimSpace(c.Phone),
		Name:  strings.TrimSpace(c.Name),
	}
}

func applyProfileHashes(userId uint64, hint capiProfileHint, out *capiUserData) {
	if out == nil || userId == 0 {
		return
	}
	if em := resolveMetaEmail(userId, hint.Email); em != "" {
		if h := hashMetaSHA256(em); h != "" {
			out.Em = []string{h}
		}
	}
	if ph := resolveMetaPhone(userId, hint.Phone); ph != "" {
		if h := hashMetaSHA256(ph); h != "" {
			out.Ph = []string{h}
		}
	}
	if fn := resolveMetaFirstName(userId, hint.Name); fn != "" {
		if h := hashMetaSHA256(fn); h != "" {
			out.Fn = []string{h}
		}
	}
	if country := resolveMetaCountry(userId); country != "" {
		if h := hashMetaSHA256(country); h != "" {
			out.Country = []string{h}
		}
	}
}

func resolveMetaEmail(userId uint64, hintEmail string) string {
	if em := normalizeMetaEmail(hintEmail); em != "" {
		return em
	}
	if ext := userinfodao.GetUserExtByUserId(userId); ext != nil {
		if em := normalizeMetaEmail(ext.Email); em != "" {
			return em
		}
	}
	acc := accountdao.GetAccountById(userId)
	if acc == nil {
		return ""
	}
	if acc.Channel == metaEmailAccountChannel {
		return normalizeMetaEmail(acc.OpenId)
	}
	return ""
}

func resolveMetaPhone(userId uint64, hintPhone string) string {
	if digits := normalizeMetaPhone("", hintPhone); digits != "" {
		return digits
	}
	info := userinfodao.GetUserInfoByUserId(userId)
	if info == nil {
		return ""
	}
	phone := strings.TrimSpace(info.Phone)
	if phone == "" {
		return ""
	}
	acc := accountdao.GetAccountById(userId)
	area := ""
	if acc != nil {
		area = strings.TrimSpace(acc.PhoneAreaCode)
	}
	return normalizeMetaPhone(area, phone)
}

func resolveMetaFirstName(userId uint64, hintName string) string {
	if fn := normalizeMetaFirstName(hintName); fn != "" {
		return fn
	}
	info := userinfodao.GetUserInfoByUserId(userId)
	if info == nil {
		return ""
	}
	return normalizeMetaFirstName(info.Nickname)
}

func resolveMetaCountry(userId uint64) string {
	loc := userloginlocationdao.GetByUserId(userId)
	if loc == nil {
		return ""
	}
	code := strings.ToLower(strings.TrimSpace(loc.LoginCountry))
	if code == "" {
		code = strings.ToLower(strings.TrimSpace(loc.RegisterCountry))
	}
	if len(code) != 2 {
		return ""
	}
	return code
}

func normalizeMetaEmail(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func normalizeMetaPhone(areaCode, phone string) string {
	digits := metaNonDigit.ReplaceAllString(phone, "")
	if digits == "" {
		return ""
	}
	areaDigits := metaNonDigit.ReplaceAllString(areaCode, "")
	if areaDigits != "" && !strings.HasPrefix(digits, areaDigits) {
		digits = areaDigits + digits
	}
	if len(digits) < 7 {
		return ""
	}
	return digits
}

func normalizeMetaFirstName(nickname string) string {
	nickname = strings.TrimSpace(nickname)
	if nickname == "" {
		return ""
	}
	if i := strings.IndexAny(nickname, " \t"); i > 0 {
		nickname = nickname[:i]
	}
	var b strings.Builder
	for _, r := range strings.ToLower(nickname) {
		if unicode.IsLetter(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func formatCAPIUserDataTags(ud capiUserData) string {
	tags := make([]string, 0, 10)
	if len(ud.ExternalID) > 0 {
		tags = append(tags, "external_id")
	}
	if len(ud.Em) > 0 {
		tags = append(tags, "em")
	}
	if len(ud.Ph) > 0 {
		tags = append(tags, "ph")
	}
	if len(ud.Fn) > 0 {
		tags = append(tags, "fn")
	}
	if len(ud.Country) > 0 {
		tags = append(tags, "country")
	}
	if ud.ClientIPAddress != "" {
		tags = append(tags, "client_ip")
	}
	if ud.ClientUserAgent != "" {
		tags = append(tags, "ua")
	}
	if ud.Fbp != "" {
		tags = append(tags, "fbp")
	}
	if ud.Fbc != "" {
		tags = append(tags, "fbc")
	}
	if len(tags) == 0 {
		return "none"
	}
	return strings.Join(tags, "+")
}
