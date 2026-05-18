package errercode

import "xr-game-server/constants/lang"

// codeTextMap 错误码的多语言文案表
// 新增/调整错误码时,需要同步在每个语言映射中维护对应文字
var codeTextMap = map[lang.Lang]map[XRCode]string{
	lang.LangEN: {
		Success:              "OK",
		EmptyUserId:          "User id is required",
		EmptyToken:           "Token is required",
		Token:                "Invalid or expired token",
		TestEnvClose:         "Test environment is closed",
		SysError:             "System error",
		LoginFail:            "Login failed",
		Ban:                  "Account banned",
		CMSLoginFail:         "CMS login failed",
		ServerClose:          "Server is closing, please try later",
		GuildExist:           "Guild name already exists",
		GuildNonExist:        "Guild does not exist",
		NoPermission:         "Permission denied",
		GuildKickSelf:        "You cannot kick yourself",
		GuildApplyExist:      "Application already exists or state is invalid",
		DiamondAmountInvalid: "Diamond amount must be positive",
		DiamondNotEnough:     "Insufficient diamond balance",
		GoldAmountInvalid:    "Gold amount must be positive",
		GoldNotEnough:        "Insufficient gold balance",
	},
	lang.LangZHCN: {
		Success:              "成功",
		EmptyUserId:          "用户ID不能为空",
		EmptyToken:           "Token不能为空",
		Token:                "Token无效或已过期",
		TestEnvClose:         "测试环境已关闭",
		SysError:             "系统错误",
		LoginFail:            "登录失败",
		Ban:                  "账号已被封禁",
		CMSLoginFail:         "后台登录失败",
		ServerClose:          "服务器关闭中,请稍后再试",
		GuildExist:           "工会名称已存在",
		GuildNonExist:        "工会不存在",
		NoPermission:         "没有权限",
		GuildKickSelf:        "不能剔除自己",
		GuildApplyExist:      "申请已存在或状态不合法",
		DiamondAmountInvalid: "钻石数量必须为正数",
		DiamondNotEnough:     "钻石余额不足",
		GoldAmountInvalid:    "金币数量必须为正数",
		GoldNotEnough:        "金币余额不足",
	},
	lang.LangZHTW: {
		Success:              "成功",
		EmptyUserId:          "用戶ID不能為空",
		EmptyToken:           "Token不能為空",
		Token:                "Token無效或已過期",
		TestEnvClose:         "測試環境已關閉",
		SysError:             "系統錯誤",
		LoginFail:            "登入失敗",
		Ban:                  "帳號已被封鎖",
		CMSLoginFail:         "後台登入失敗",
		ServerClose:          "伺服器關閉中,請稍後再試",
		GuildExist:           "公會名稱已存在",
		GuildNonExist:        "公會不存在",
		NoPermission:         "沒有權限",
		GuildKickSelf:        "不能剔除自己",
		GuildApplyExist:      "申請已存在或狀態不合法",
		DiamondAmountInvalid: "鑽石數量必須為正數",
		DiamondNotEnough:     "鑽石餘額不足",
		GoldAmountInvalid:    "金幣數量必須為正數",
		GoldNotEnough:        "金幣餘額不足",
	},
}

// GetMsg 返回错误码在指定语言下的提示文字;
// 未匹配目标语言则回落到默认语言;再未匹配则返回空字符串
func GetMsg(code XRCode, l lang.Lang) string {
	if m, ok := codeTextMap[l]; ok {
		if s, ok2 := m[code]; ok2 {
			return s
		}
	}
	if m, ok := codeTextMap[lang.DefaultLang]; ok {
		if s, ok2 := m[code]; ok2 {
			return s
		}
	}
	return ""
}
