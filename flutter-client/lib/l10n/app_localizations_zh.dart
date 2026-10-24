// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for Chinese (`zh`).
class AppLocalizationsZh extends AppLocalizations {
  AppLocalizationsZh([String locale = 'zh']) : super(locale);

  @override
  String get appTitle => 'Sara Live';

  @override
  String get uiRenderErrorTitle => 'UI 渲染错误';

  @override
  String get commonCancel => '取消';

  @override
  String get commonSave => '保存';

  @override
  String get commonRetry => '重试';

  @override
  String get commonCopied => '已复制';

  @override
  String get inviteCodeCopied => '邀请码已复制';

  @override
  String get languageTitle => '语言';

  @override
  String get languageSubtitle => '中文 / English';

  @override
  String get languageEnglish => 'English';

  @override
  String get languageChinese => '中文';

  @override
  String get languagePickerTitle => '语言';

  @override
  String get languagePickerHint => '选择后自动保存';

  @override
  String get uiStyleTitle => '界面风格';

  @override
  String get uiStyleSubtitle => '点击切换主题';

  @override
  String get settingsSection => '设置';

  @override
  String get accountSection => '账户';

  @override
  String get logout => '退出登录';

  @override
  String get editNicknameTitle => '修改昵称';

  @override
  String get editNicknameHint => '请输入新昵称';

  @override
  String get avatarFormatError => '仅支持 jpg/png/webp/gif/bmp';

  @override
  String userDisplayName(int userId) {
    return '用户$userId';
  }

  @override
  String get walletGold => '金币';

  @override
  String get walletDiamond => '钻石';

  @override
  String get walletFollow => '关注';

  @override
  String get walletFans => '粉丝';

  @override
  String get vipTitle => 'VIP';

  @override
  String get inviteCodeTitle => '邀请码';

  @override
  String get anchorRoleTitle => '主播身份';

  @override
  String get anchorRoleEnabled => '已开通';

  @override
  String get anchorRoleDisabled => '未开通';

  @override
  String get tabOneToOne => '1v1';

  @override
  String get tabLive => '直播';

  @override
  String get tabGame => '游戏';

  @override
  String get tabMine => '我的';

  @override
  String get tabOnlineAudience => '1v1';

  @override
  String get tabOneToOneRoom => '1v1房';

  @override
  String get tabShowLive => '直播';

  @override
  String get placeholderLiveTitle => '直播';

  @override
  String get placeholderLiveSubtitle => '热门秀场与游戏直播列表稍后接入';

  @override
  String get placeholderGameTitle => '游戏';

  @override
  String get placeholderGameSubtitle => '游戏货架与半屏玩法稍后接入';

  @override
  String get deviceCodeTitle => '设备码';

  @override
  String get deviceCodeEditTitle => '编辑设备码';

  @override
  String get deviceCodeHint => '用于设备码快捷登录（最多 128 字符）';

  @override
  String get deviceCodeWarning => '保存后将使用该设备码重新登录，可能切换到其他账号';

  @override
  String get deviceCodeEmpty => '设备码不能为空';

  @override
  String get deviceCodeTooLong => '设备码不能超过 128 个字符';

  @override
  String get deviceCodeUnchanged => '设备码未变更';

  @override
  String get deviceCodeSwitchSuccess => '已切换设备码并登录';

  @override
  String get deviceCodeNotSet => '未设置（登录时自动生成）';

  @override
  String get apiBaseTitle => 'API 服务器';

  @override
  String get apiBaseEditTitle => '编辑 API 地址';

  @override
  String get apiBaseHint => '例如 http://192.168.2.11:9443';

  @override
  String get apiBaseWarning => '保存后会切换服务器并退出登录，需在新地址上重新登录。';

  @override
  String get apiBaseResetDefault => '恢复默认';

  @override
  String get apiBaseEmpty => 'API 地址不能为空';

  @override
  String get apiBaseInvalid => '请输入有效的 http 或 https 地址';

  @override
  String get apiBaseUnchanged => 'API 地址未变更';

  @override
  String get apiBaseSwitchSuccess => '已切换 API 服务器，请重新登录';

  @override
  String get apiBaseNotSet => '未配置';

  @override
  String apiBaseUsingDefault(String url) {
    return '默认：$url';
  }

  @override
  String get liveHotListTitle => '正在直播';

  @override
  String get liveHotListEmpty => '暂无直播中的房间';

  @override
  String get liveBadgeOnAir => '直播中';

  @override
  String liveOnlineCount(int count) {
    return '$count 人在看';
  }

  @override
  String get liveChatHint => '说点什么…';

  @override
  String get liveSendGift => '送礼物';

  @override
  String get liveGiftListEmpty => '暂无礼物';

  @override
  String get liveAnchorEnded => '主播已下播';

  @override
  String get liveWatchBlockedInCall => '请先结束通话再观看直播';

  @override
  String get liveJoiningRoom => '正在进入直播间…';

  @override
  String get liveWatchAudienceTitle => '在线观众';

  @override
  String liveWatchAudienceSubtitle(int count) {
    return '当前 $count 人在看（不含主播）';
  }

  @override
  String get liveWatchAudienceEmpty => '暂无在线观众';

  @override
  String get liveWatchTapToChat => '聊天';

  @override
  String liveGiftSent(String giftName, double diamond) {
    return '已送出 $giftName，余额 $diamond';
  }
}
