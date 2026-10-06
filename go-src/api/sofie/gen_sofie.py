# -*- coding: utf-8 -*-
"""Generate Sofie App route wrappers and README. Run from repo root: python go-src/api/sofie/gen_sofie.py"""
from collections import defaultdict
from pathlib import Path

ROOT = Path(__file__).resolve().parent

# old_url, auth(yes/no), controller_expr, summary, sofie_module, sofie_action
APIS = [
    ("/auth/sendCode", "no", "(&controller.AuthController{}).SendCode", "发送短信验证码", "access", "dispatchSmsOtp"),
    ("/auth/sendEmailCode", "no", "(&controller.AuthController{}).SendEmailCode", "发送邮箱验证码", "access", "dispatchMailOtp"),
    ("/auth/emailLogin", "no", "(&controller.AuthController{}).EmailLogin", "邮箱登录", "access", "signInByMailbox"),
    ("/auth/firebaseLogin", "no", "(&controller.AuthController{}).FirebaseLogin", "Firebase登录", "access", "signInByFirebase"),
    ("/auth/phoneRegister", "no", "(&controller.AuthController{}).PhoneRegister", "手机号注册", "access", "createPhoneAccount"),
    ("/auth/phoneLogin", "no", "(&controller.AuthController{}).PhoneLogin", "手机号登录", "access", "signInByMobile"),
    ("/auth/coinMerchantLogin", "no", "(&controller.AuthController{}).CoinMerchantLogin", "币商登录", "access", "signInByDealer"),
    ("/auth/deviceLogin", "no", "(&controller.AuthController{}).DeviceLogin", "设备码登录", "access", "signInByDeviceCode"),
    ("/auth/h5DeviceLogin", "no", "(&controller.AuthController{}).H5DeviceLogin", "H5设备码登录", "access", "signInByH5Device"),
    ("/auth/phoneResetPassword", "no", "(&controller.AuthController{}).PhoneResetPassword", "手机号重置密码", "access", "resetMobileSecret"),
    ("/auth/phoneChangePassword", "yes", "(&controller.AuthAppController{}).PhoneChangePassword", "修改登录密码", "access", "replaceMobileSecret"),
    ("/auth/bindEmail", "yes", "(&controller.AuthAppController{}).BindEmail", "绑定邮箱", "access", "linkMailbox"),
    ("/auth/bindFirebase", "yes", "(&controller.AuthAppController{}).BindFirebase", "绑定Firebase", "access", "linkFirebase"),
    ("/userInfo/get", "yes", "(&controller.UserInfoController{}).Get", "获取当前用户信息", "profile", "fetchMine"),
    ("/userInfo/getUserExt", "yes", "(&controller.UserInfoController{}).GetUserExt", "获取用户扩展信息", "profile", "fetchExtra"),
    ("/userInfo/updateNickname", "yes", "(&controller.UserInfoController{}).UpdateNickname", "修改昵称", "profile", "renameNick"),
    ("/userInfo/updateGender", "yes", "(&controller.UserInfoController{}).UpdateGender", "修改性别", "profile", "changeSex"),
    ("/userInfo/updateBirthday", "yes", "(&controller.UserInfoController{}).UpdateBirthday", "修改生日", "profile", "changeBirthDate"),
    ("/userInfo/getCurrencyLog", "yes", "(&controller.UserInfoController{}).GetCurrencyLog", "查询货币流水", "profile", "browseLedger"),
    ("/userInfo/cancelAccount", "yes", "(&controller.UserInfoController{}).CancelAccount", "注销账号", "profile", "closeAccount"),
    ("/userInfo/feedback", "yes", "(&controller.UserInfoController{}).Feedback", "提交反馈", "profile", "submitOpinion"),
    ("/userInfo/report", "yes", "(&controller.UserInfoController{}).Report", "举报用户", "profile", "submitComplaint"),
    ("/userInfo/reportAttribution", "yes", "(&controller.UserInfoController{}).ReportAttribution", "上报归因", "profile", "submitAttribution"),
    ("/userInfo/getByInviteCode", "yes", "(&controller.UserInfoController{}).GetUserInfoByInviteCode", "按邀请码查用户", "profile", "lookupByInvite"),
    ("/userInfo/reportInviter", "yes", "(&controller.UserInfoController{}).ReportInviter", "上报邀请人", "profile", "submitInviter"),
    ("/userInfo/cancelAccountByCode", "no", "(&controller.UserInfoPublicController{}).CancelAccountByCode", "官网验证码销户", "profile", "closeAccountByOtp"),
    ("/userInfo/uploadAvatar", "yes", "RAW:controller.HandleUploadAvatar", "上传头像", "profile", "changePhoto"),
    ("/liveRoom/create", "yes", "RAW:controller.HandleCreateLiveRoom", "创建/更新直播间", "studio", "openBooth"),
    ("/liveRoom/createOneToOne", "yes", "(&controller.LiveRoomAppController{}).CreateOneToOneRoom", "开通1v1房间", "studio", "openDirectBooth"),
    ("/liveRoom/getOneToOne", "yes", "(&controller.LiveRoomAppController{}).GetOneToOneRoom", "查询1v1房间配置", "studio", "fetchDirectBooth"),
    ("/liveRoom/startLive", "yes", "(&controller.LiveRoomAppController{}).StartLive", "开播", "studio", "beginBroadcast"),
    ("/liveRoom/stopLive", "yes", "(&controller.LiveRoomAppController{}).StopLive", "下播", "studio", "endBroadcast"),
    ("/liveRoom/updateCover", "yes", "(&controller.LiveRoomAppController{}).UpdateCover", "修改封面", "studio", "replaceCover"),
    ("/liveRoom/updateNotice", "yes", "(&controller.LiveRoomAppController{}).UpdateNotice", "修改公告", "studio", "replaceNotice"),
    ("/liveRoom/join", "yes", "(&controller.LiveRoomAppController{}).JoinRoom", "加入直播间", "studio", "enterHall"),
    ("/liveRoom/leave", "yes", "(&controller.LiveRoomAppController{}).LeaveRoom", "离开直播间", "studio", "exitHall"),
    ("/liveRoom/getLiveRecord", "yes", "(&controller.LiveRoomAppController{}).GetLiveRecord", "查询本场直播记录", "studio", "fetchRound"),
    ("/liveRoom/onlineList", "yes", "(&controller.LiveRoomAppController{}).GetOnlineUserList", "在线观众列表", "studio", "browseAudience"),
    ("/liveRoom/vipOnlineList", "yes", "(&controller.LiveRoomAppController{}).GetVipOnlineUserList", "VIP在线观众", "studio", "browseVipAudience"),
    ("/liveRoom/contributionRank", "yes", "(&controller.LiveRoomAppController{}).GetContributionRank", "贡献榜", "studio", "browseContribution"),
    ("/liveRoom/get", "yes", "(&controller.LiveRoomAppController{}).GetRoom", "查询直播间", "studio", "fetchBooth"),
    ("/liveRoom/roomList", "yes", "(&controller.LiveRoomAppController{}).RoomList", "直播间列表", "studio", "browseBooths"),
    ("/liveRoom/oneToOneRoomList", "yes", "(&controller.LiveRoomAppController{}).OneToOneRoomList", "1v1房间列表", "studio", "browseDirectBooths"),
    ("/liveRoom/serverOnlineNormalUserList", "yes", "(&controller.LiveRoomAppController{}).ServerOnlineNormalUserList", "服务器在线普通用户", "studio", "browseOnlineNormals"),
    ("/liveRoom/nearbyRoomList", "yes", "(&controller.LiveRoomAppController{}).NearbyRoomList", "相邻直播间", "studio", "browseNearbyBooths"),
    ("/liveRoom/hotRoomList", "yes", "(&controller.LiveRoomAppController{}).HotRoomList", "Hot房间列表", "studio", "browseHotBooths"),
    ("/liveRoom/followedRoomList", "yes", "(&controller.LiveRoomAppController{}).FollowedRoomList", "关注的直播间", "studio", "browseFollowedBooths"),
    ("/liveRoom/sendGift", "yes", "(&controller.LiveRoomAppController{}).SendGift", "直播间送礼", "studio", "postPresent"),
    ("/liveRoom/sendGiftToAnchor", "yes", "(&controller.LiveRoomAppController{}).SendGiftToAnchor", "给指定主播送礼", "studio", "postPresentToHost"),
    ("/liveRoom/sendChat", "yes", "(&controller.LiveRoomAppController{}).SendChat", "直播间文字消息", "studio", "postChat"),
    ("/liveRoom/sendPaidDanmaku", "yes", "(&controller.LiveRoomAppController{}).SendPaidDanmaku", "付费弹幕", "studio", "postPaidBullet"),
    ("/liveRoom/setAudienceMute", "yes", "(&controller.LiveRoomAppController{}).SetAudienceMute", "禁言观众", "studio", "silenceViewer"),
    ("/liveRoom/cancelAudienceMute", "yes", "(&controller.LiveRoomAppController{}).CancelAudienceMute", "取消禁言", "studio", "unsilenceViewer"),
    ("/liveRoom/getAudienceRestrictStatus", "yes", "(&controller.LiveRoomAppController{}).GetAudienceRestrictStatus", "查询观众限制状态", "studio", "fetchViewerRestrict"),
    ("/liveRoom/kickAudience", "yes", "(&controller.LiveRoomAppController{}).KickAudience", "踢出观众", "studio", "removeViewer"),
    ("/liveRoom/cancelKickBan", "yes", "(&controller.LiveRoomAppController{}).CancelKickBan", "取消进入限制", "studio", "clearViewerBan"),
    ("/liveRoom/reportLiveStartStatus", "yes", "(&controller.LiveRoomAppController{}).ReportLiveStartStatus", "上报开播状态", "studio", "reportBroadcastState"),
    ("/liveRoom/gameRecommendList", "yes", "(&controller.LiveRoomAppController{}).GameRecommendList", "推荐游戏列表", "studio", "browseGamePicks"),
    ("/liveRoom/reportAnchorCode", "yes", "(&controller.LiveRoomAppController{}).ReportAnchorCode", "上报主播码", "studio", "submitHostCode"),
    ("/call/liveRoomCall", "yes", "(&controller.CallAppController{}).LiveRoomCall", "直播间发起通话", "session", "inviteFromHall"),
    ("/call/oneToOneRoomCall", "yes", "(&controller.CallAppController{}).OneToOneRoomCall", "1v1房间发起通话", "session", "inviteDirectVideo"),
    ("/call/batchInviteLiveRoomCall", "yes", "(&controller.CallAppController{}).BatchInviteLiveRoomCall", "批量邀请通话", "session", "inviteManyFromHall"),
    ("/call/anchorRejectCall", "yes", "(&controller.CallAppController{}).AnchorRejectCall", "拒接通话", "session", "declineInvite"),
    ("/call/acceptCall", "yes", "(&controller.CallAppController{}).AcceptCall", "接听通话", "session", "acceptInvite"),
    ("/call/confirmCall", "yes", "(&controller.CallAppController{}).ConfirmCall", "通话应答确认", "session", "ackInvite"),
    ("/call/getCallConfirmStatus", "yes", "(&controller.CallAppController{}).GetCallConfirmStatus", "查询应答确认状态", "session", "fetchAckState"),
    ("/call/callTimeout", "yes", "(&controller.CallAppController{}).CallTimeout", "呼叫超时", "session", "markInviteTimeout"),
    ("/call/callHeart", "yes", "(&controller.CallAppController{}).CallHeart", "通话心跳", "session", "keepAlive"),
    ("/call/endCall", "yes", "(&controller.CallAppController{}).EndCall", "结束通话", "session", "hangUp"),
    ("/message/sendPrivateMessage", "yes", "(&controller.MessageAppController{}).SendPrivateMessage", "发送私信", "inbox", "postDirectNote"),
    ("/message/privateMessageUnreadList", "yes", "(&controller.MessageAppController{}).PrivateMessageUnreadList", "私信未读列表", "inbox", "browseUnreadNotes"),
    ("/message/privateMessageBySender", "yes", "(&controller.MessageAppController{}).PrivateMessageBySender", "按发送者查私信", "inbox", "browseNotesByPeer"),
    ("/message/clearPrivateMessageUnread", "yes", "(&controller.MessageAppController{}).ClearPrivateMessageUnread", "清除指定私信未读", "inbox", "markPeerNotesRead"),
    ("/message/clearAllPrivateMessageUnread", "yes", "(&controller.MessageAppController{}).ClearAllPrivateMessageUnread", "清除全部私信未读", "inbox", "markAllNotesRead"),
    ("/message/batchDeletePrivateMessage", "yes", "(&controller.MessageAppController{}).BatchDeletePrivateMessage", "批量删除私信", "inbox", "removeNotesBatch"),
    ("/message/personalSystemMessageList", "yes", "(&controller.MessageAppController{}).PersonalSystemMessageList", "个人系统消息", "inbox", "browsePersonalNotices"),
    ("/message/activityMessageList", "yes", "(&controller.MessageAppController{}).ActivityMessageList", "活动消息列表", "inbox", "browseCampaignNotices"),
    ("/message/systemMessageUnreadList", "yes", "(&controller.MessageAppController{}).SystemMessageUnreadList", "系统消息未读", "inbox", "browseUnreadNotices"),
    ("/message/clearSystemMessageUnread", "yes", "(&controller.MessageAppController{}).ClearSystemMessageUnread", "系统消息已读", "inbox", "markNoticesRead"),
    ("/message/messageUnreadCount", "yes", "(&controller.MessageAppController{}).MessageUnreadCount", "未读数", "inbox", "countUnread"),
    ("/shortVideo/appPublishShortVideo", "yes", "RAW:controller.HandleAppPublishShortVideo", "发布短视频", "clip", "submitMedia"),
    ("/shortVideo/appShortVideoList", "yes", "(&controller.ShortVideoAppController{}).AppShortVideoList", "短视频列表", "clip", "browseFeed"),
    ("/shortVideo/appShortVideoScroll", "yes", "(&controller.ShortVideoAppController{}).AppShortVideoScroll", "短视频滑动列表", "clip", "browseFeedScroll"),
    ("/shortVideo/appShortVideoViewList", "yes", "(&controller.ShortVideoAppController{}).AppShortVideoViewList", "短视频观看列表", "clip", "browseViewed"),
    ("/shortVideo/appShortVideoPublishList", "yes", "(&controller.ShortVideoAppController{}).AppShortVideoPublishList", "我发布的短视频", "clip", "browsePublished"),
    ("/shortVideo/likeShortVideo", "yes", "(&controller.ShortVideoAppController{}).LikeShortVideo", "点赞短视频", "clip", "toggleLike"),
    ("/shortVideo/appShortVideoCfg", "yes", "(&controller.ShortVideoAppController{}).AppShortVideoCfg", "短视频配置", "clip", "fetchSettings"),
    ("/shortVideo/watchShortVideoStart", "yes", "(&controller.ShortVideoAppController{}).WatchShortVideoStart", "开始观看", "clip", "beginWatch"),
    ("/shortVideo/watchShortVideoEnd", "yes", "(&controller.ShortVideoAppController{}).WatchShortVideoEnd", "结束观看", "clip", "finishWatch"),
    ("/shortVideo/payShortVideo", "yes", "(&controller.ShortVideoAppController{}).PayShortVideo", "付费观看", "clip", "unlockMedia"),
    ("/shortVideo/appShortVideoWatchList", "yes", "(&controller.ShortVideoAppController{}).AppShortVideoWatchList", "观看记录", "clip", "browseWatchLog"),
    ("/shortVideo/appShortVideoCategoryList", "yes", "(&controller.ShortVideoAppController{}).AppShortVideoCategoryList", "短视频分类", "clip", "browseCategories"),
    ("/shortVideo/appShortVideoUploadRecordList", "yes", "(&controller.ShortVideoAppController{}).AppShortVideoUploadRecordList", "上传记录", "clip", "browseUploads"),
    ("/shortVideo/appShortVideoPendingReviewList", "yes", "(&controller.ShortVideoAppController{}).AppShortVideoPendingReviewList", "待审列表", "clip", "browsePendingReview"),
    ("/shortVideo/appDeleteShortVideo", "yes", "(&controller.ShortVideoAppController{}).AppDeleteShortVideo", "删除短视频", "clip", "removeMedia"),
    ("/shortVideo/appShortVideoStatList", "yes", "(&controller.ShortVideoAppController{}).AppShortVideoStatList", "短视频统计", "clip", "browseStats"),
    ("/shortVideo/appShortVideoPriceTierList", "yes", "(&controller.ShortVideoAppController{}).AppShortVideoPriceTierList", "价格档位", "clip", "browsePriceTiers"),
    ("/shortVideo/appUserShortVideoList", "yes", "(&controller.ShortVideoAppController{}).AppUserShortVideoList", "指定用户短视频", "clip", "browseUserFeed"),
    ("/game/appGameList", "yes", "(&controller.GameAppController{}).AppGameList", "游戏列表", "play", "browseTitles"),
    ("/game/appGameStart", "yes", "(&controller.GameAppController{}).AppGameStart", "启动游戏", "play", "openTitle"),
    ("/game/appGameBetList", "yes", "(&controller.GameAppController{}).AppGameBetList", "下注记录", "play", "browseBets"),
    ("/game/appGameConsumeRank", "yes", "(&controller.GameAppController{}).AppGameConsumeRank", "单场游戏消费榜", "play", "browseRoundSpend"),
    ("/gift/appGiftList", "yes", "(&controller.GiftAppController{}).GiftList", "礼物列表", "present", "browseCatalog"),
    ("/rechargeOrder/createRechargeOrder", "yes", "(&controller.RechargeOrderAppController{}).CreateRechargeOrder", "创建充值订单", "topup", "openStoreOrder"),
    ("/rechargeOrder/createChannelRechargeOrder", "yes", "(&controller.RechargeOrderAppController{}).CreateChannelRechargeOrder", "渠道充值建单", "topup", "openChannelOrder"),
    ("/rechargeOrder/createCoinMerchantChannelRechargeOrder", "yes", "(&controller.RechargeOrderAppController{}).CreateCoinMerchantChannelRechargeOrder", "币商渠道建单", "topup", "openDealerChannelOrder"),
    ("/rechargeOrder/getChannelPayUserProfile", "yes", "(&controller.RechargeOrderAppController{}).GetChannelPayUserProfile", "查询付款人资料", "topup", "fetchPayerProfile"),
    ("/rechargeOrder/saveChannelPayUserProfile", "yes", "(&controller.RechargeOrderAppController{}).SaveChannelPayUserProfile", "保存付款人资料", "topup", "savePayerProfile"),
    ("/rechargeOrder/myRechargeOrderList", "yes", "(&controller.RechargeOrderAppController{}).MyRechargeOrderList", "我的充值订单", "topup", "browseMyOrders"),
    ("/rechargeOrder/checkRechargeOrderSuccess", "yes", "(&controller.RechargeOrderAppController{}).CheckRechargeOrderSuccess", "查询充值是否成功", "topup", "checkOrderPaid"),
    ("/rechargeCfg/rechargeCfgListForApp", "yes", "(&controller.RechargeCfgAppController{}).RechargeCfgList", "充值配置列表", "catalog", "browsePackages"),
    ("/rechargeCfg/rechargeCfgListByUserId", "no", "(&controller.RechargeCfgAppPublicController{}).RechargeCfgListByUserId", "按用户ID查充值配置", "catalog", "browsePackagesByUser"),
    ("/ticket/appTicketList", "yes", "(&controller.TicketAppController{}).AppTicketList", "门票列表", "pass", "browseTickets"),
    ("/privateRoomBilling/appBillingList", "yes", "(&controller.PrivateRoomBillingAppController{}).AppBillingList", "1v1计费档位", "rate", "browseMinutePrices"),
    ("/vip/getVipDetail", "yes", "(&controller.VipAppController{}).GetVipDetail", "VIP详情", "member", "fetchStatus"),
    ("/vipCfg/getVipCfgByLevel", "yes", "(&controller.VipCfgAppController{}).GetVipCfgByLevel", "按等级查VIP配置", "memberPlan", "fetchTier"),
    ("/vipCfg/vipCfgListForApp", "yes", "(&controller.VipCfgAppController{}).VipCfgListForApp", "VIP配置列表", "memberPlan", "browseTiers"),
    ("/banner/appBannerList", "yes", "(&controller.BannerAppController{}).AppBannerList", "首页Banner", "poster", "browseHome"),
    ("/agora/liveRoomToken", "yes", "(&controller.AgoraAppController{}).LiveRoomToken", "声网房间Token", "rtc", "issueChannelToken"),
    ("/agora/appId", "yes", "(&controller.AgoraAppController{}).AppId", "声网AppId", "rtc", "fetchVendorId"),
    ("/liveRoomTag/appLiveRoomTagList", "yes", "(&controller.LiveRoomTagAppController{}).AppLiveRoomTagList", "直播间标签", "label", "browseAll"),
    ("/liveRoomTag/appLiveRoomNormalTagList", "yes", "(&controller.LiveRoomTagAppController{}).AppLiveRoomNormalTagList", "普通直播间标签", "label", "browseNormal"),
    ("/liveFollow/follow", "yes", "(&controller.LiveFollowAppController{}).Follow", "关注主播", "social", "addHost"),
    ("/liveFollow/unfollow", "yes", "(&controller.LiveFollowAppController{}).Unfollow", "取消关注", "social", "removeHost"),
    ("/liveFollow/isFollowing", "yes", "(&controller.LiveFollowAppController{}).IsFollowing", "是否已关注", "social", "checkHost"),
    ("/liveFollow/followingList", "yes", "(&controller.LiveFollowAppController{}).FollowingList", "关注列表", "social", "browseHosts"),
    ("/liveFollow/followerList", "yes", "(&controller.LiveFollowAppController{}).FollowerList", "粉丝列表", "social", "browseFans"),
    ("/liveFollow/block", "yes", "(&controller.LiveFollowAppController{}).Block", "拉黑用户", "social", "banUser"),
    ("/liveFollow/unblock", "yes", "(&controller.LiveFollowAppController{}).Unblock", "取消拉黑", "social", "unbanUser"),
    ("/liveFollow/blockList", "yes", "(&controller.LiveFollowAppController{}).BlockList", "拉黑列表", "social", "browseBanned"),
    ("/liveRecord/appLiveRecordList", "yes", "(&controller.LiveRecordAppController{}).AppLiveRecordList", "主播直播记录", "history", "browseBroadcasts"),
    ("/richRank/appRichRankList", "yes", "(&controller.RichRankAppController{}).AppRichRankList", "富豪榜", "wealth", "browseBoard"),
    ("/anchorRank/appAnchorRankList", "yes", "(&controller.AnchorRankAppController{}).AppAnchorRankList", "主播红人榜", "star", "browseBoard"),
    ("/gameConsumeRank/appGameConsumeRankList", "yes", "(&controller.GameConsumeRankAppController{}).AppGameConsumeRankList", "游戏消费榜", "spend", "browseBoard"),
    ("/gold/exchangeGoldToDiamond", "yes", "(&controller.GoldAppController{}).ExchangeGoldToDiamond", "金币兑换钻石", "coin", "swapToGem"),
    ("/gold/transferGold", "yes", "(&controller.GoldAppController{}).TransferGold", "转赠金币", "coin", "sendToUser"),
    ("/gold/getCoinMerchantTransferRecordList", "yes", "(&controller.GoldAppController{}).GetCoinMerchantTransferRecordList", "币商转账记录", "coin", "browseDealerTransfers"),
    ("/firstRechargeActivity/firstRechargeActivityCfgForApp", "yes", "(&controller.FirstRechargeActivityAppController{}).FirstRechargeActivityCfgForApp", "首充活动配置", "bonus", "fetchFirstPay"),
    ("/appVersion/appVersionQuery", "no", "(&controller.AppVersionAppController{}).AppVersionQuery", "App版本查询", "release", "lookupBuild"),
    ("/customerService/cfg", "no", "(&controller.CustomerServiceAppController{}).Cfg", "客服配置", "help", "fetchContact"),
    ("/fiatCurrency/fiatCurrencyListForApp", "no", "(&controller.FiatCurrencyAppController{}).FiatCurrencyListForApp", "法币列表", "money", "browseFiats"),
    ("/firebase/getClientCfgForApp", "no", "(&controller.FirebaseAppController{}).GetClientCfgForApp", "Firebase客户端配置", "googleAuth", "fetchClientCfg"),
    ("/coinMerchantRechargeCfg/coinMerchantRechargeCfgListForApp", "yes", "(&controller.CoinMerchantRechargeCfgAppController{}).CoinMerchantRechargeCfgList", "币商充值档位", "dealer", "browsePackages"),
    ("/coinMerchantRechargeCfg/paymentRegionList", "yes", "(&controller.CoinMerchantRechargeCfgAppController{}).PaymentRegionList", "币商支付区域", "dealer", "browsePayRegions"),
    ("/sysInfo/cfg", "no", "(&controller.SysInfoController{}).GetInfo", "系统配置", "boot", "loadConfig"),
]


def validate():
    olds, news, last_new, last_old = set(), set(), set(), set()
    for old, _auth, _fn, _sum, mod, act in APIS:
        if old in olds:
            raise SystemExit(f"duplicate old {old}")
        olds.add(old)
        new = f"/sofie/{mod}/{act}"
        if new in news:
            raise SystemExit(f"duplicate new {new}")
        news.add(new)
        old_last = old.rsplit("/", 1)[-1]
        if act == old_last:
            raise SystemExit(f"action same as old last: {old} -> {act}")
        if mod == old.strip("/").split("/")[0]:
            raise SystemExit(f"module same as old prefix: {old} -> {mod}")
        last_new.add(act)
        last_old.add(old_last)
    overlap = last_new & last_old
    if overlap:
        raise SystemExit(f"new action collides with some old last segment: {sorted(overlap)}")


def handler_expr(fn: str) -> str:
    if fn.startswith("RAW:"):
        return fn[4:]
    return f"jsonHandler({fn})"


def write_go():
    groups = defaultdict(lambda: {"yes": [], "no": []})
    for old, auth, fn, summary, mod, act in APIS:
        groups[mod][auth].append((old, fn, summary, act))

    lines = [
        "package sofie",
        "",
        "// Code generated by gen_sofie.py; DO NOT EDIT.",
        "",
        "import (",
        '\t"xr-game-server/controller"',
        '\t"xr-game-server/core/httpserver"',
        ")",
        "",
        "func Init() {",
    ]
    for mod in sorted(groups):
        for auth in ("yes", "no"):
            items = groups[mod][auth]
            if not items:
                continue
            reg = "RegAppRouteGroup" if auth == "yes" else "RegNonAuthAppRouteGroup"
            lines.append(f'\thttpserver.{reg}("/sofie/{mod}", []httpserver.AppRoute{{')
            for _old, fn, _summary, act in items:
                lines.append(f'\t\t{{Path: "/{act}", Handler: {handler_expr(fn)}}},')
            lines.append("\t})")
    lines.append("\tsetupSofieOpenApi()")
    lines.append("}")
    lines.append("")
    (ROOT / "routes_gen.go").write_text("\n".join(lines) + "\n", encoding="utf-8")


def old_http_method(old: str) -> str:
    if old == "/sysInfo/cfg":
        return "get"
    return "post"


def write_openapi_map():
    lines = [
        "package sofie",
        "",
        "// Code generated by gen_sofie.py; DO NOT EDIT.",
        "",
        "type openAPIDocRoute struct {",
        "\tOldMethod string",
        "\tOldPath   string",
        "\tNewPath   string",
        "\tTag       string",
        "\tSummary   string",
        "}",
        "",
        "var openAPIDocRoutes = []openAPIDocRoute{",
    ]
    for old, _auth, _fn, summary, mod, act in APIS:
        method = old_http_method(old)
        new = f"/sofie/{mod}/{act}"
        lines.append(
            f'\t{{OldMethod: "{method}", OldPath: "{old}", NewPath: "{new}", Tag: "{mod}", Summary: "{summary}"}},'
        )
    lines.append("}")
    lines.append("")
    (ROOT / "openapi_map_gen.go").write_text("\n".join(lines) + "\n", encoding="utf-8")


def write_md():
    rows = [
        "# Sofie App 接口对接表",
        "",
        "新旧并存：旧路径仍可用。Sofie 路径统一 `POST /sofie/{模块}/{动作}`，请求体字段与旧接口相同。",
        "",
        "| 模块 | 说明 | 鉴权 | 旧 URL | 新 URL |",
        "| --- | --- | --- | --- | --- |",
    ]
    for old, auth, _fn, summary, mod, act in APIS:
        rows.append(
            f"| `{mod}` | {summary} | {'需要登录' if auth == 'yes' else '免登录'} | `POST {old}` | `POST /sofie/{mod}/{act}` |"
        )
    rows.append("")
    rows.append("说明：")
    rows.append("")
    rows.append("- `/sysInfo/cfg` 旧接口为 GET，Sofie 对应接口改为 POST，请求体仍按原字段。")
    rows.append("- 创建房间、上传头像、发布短视频仍为 multipart，路径换成 Sofie 新地址即可。")
    rows.append("- Apifox/OpenAPI 仅 `config/local` 开启，独立地址见 `server.sofieOpenapiPath` / `server.sofieSwaggerPath`，不与旧 App 文档混用。")
    rows.append("")
    (ROOT / "README.md").write_text("\n".join(rows), encoding="utf-8")


if __name__ == "__main__":
    validate()
    write_go()
    write_openapi_map()
    write_md()
    print(f"generated {len(APIS)} sofie routes")
