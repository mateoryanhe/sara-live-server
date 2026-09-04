-- =============================================================================
-- live_db 上线前清理脚本：删除用户/业务/统计数据表，保留配置类表
-- =============================================================================
-- 生成依据：go-src/entity 下 TbName 清单（2026-09-04）
--
-- 【危险】DROP TABLE 不可逆。执行前务必：
--   1. 确认当前库名（勿误连正式未备份库）
--   2. 已做全量备份并可恢复
--   3. 停写或停服务后再执行
--   4. 将下方 SAFETY_OK 改为 1 才会真正 DROP
--
-- 执行后：重启 Go；配置表仍在，业务表由 AutoMigrate 按需重建空表。
-- CMS 账号(cms_users/cms_roles/permissions)默认保留，便于继续登录后台。
-- =============================================================================

SELECT DATABASE() AS current_database;

-- >>> 安全开关：改成 1 才执行 DROP <<<
SET @SAFETY_OK := 0;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- -----------------------------------------------------------------------------
-- 【保留】配置 / CMS 后台（本脚本不 DROP，仅作清单）
-- -----------------------------------------------------------------------------
-- account_cfgs                         账号相关开关配置
-- agora_cfgs                           声网配置
-- aliyun_text_moderation_cfgs          阿里云文本审核配置
-- app_pkgs                             App 包名/渠道包配置
-- app_version_cfgs                     App 版本更新配置
-- app_version_update_details           App 版本更新文案明细
-- cms_roles                            CMS 角色
-- cms_users                            CMS 后台账号
-- permissions                          CMS 权限点
-- customer_service_cfgs                客服配置
-- data_sync_cfgs                       跨环境数据同步配置
-- fiat_currency_cfgs                   法币币种配置
-- first_recharge_activity_cfgs         首充活动配置
-- first_recharge_activity_privileges   首充活动特权配置
-- game_cfgs                            游戏配置
-- game_platform_cfgs                   游戏平台接入配置
-- vendor_game_libs                     厂商游戏库/列表缓存配置侧
-- google_play_cfgs                     Google Play 充值配置
-- h5_live_deploy_cfgs                  H5 直播部署配置
-- home_banners                         首页 Banner 运营配置
-- activity_messages                    活动消息（运营下发模板/内容）
-- live_cfgs                            直播全局配置
-- live_gifts                           礼物配置
-- live_private_room_billings           1v1/私密计费档位配置
-- live_revenue_share_cfgs              直播分成配置
-- live_room_tags                       直播间标签配置
-- live_tickets                         门票配置
-- anchor_salary_cfgs                   主播底薪配置
-- preload_cfgs                         预加载/热更相关配置
-- privacy_policy_cfgs                  隐私政策配置
-- random_nicknames                     随机昵称词库
-- recharge_cfgs                        充值档位配置
-- short_video_categories               短视频分类配置
-- short_video_cfgs                     短视频全局配置
-- short_video_price_tiers              短视频价格档配置
-- simulator_cpu_keywords               模拟器 CPU 关键词配置
-- upload_resource_cfgs                 上传资源配置
-- vip_cfgs                             VIP 配置
-- wallet_exchange_cfgs                 钱包兑换配置
-- yhpay_cfgs                           第三方支付(YhPay)配置

-- -----------------------------------------------------------------------------
-- 【删除】用户体系
-- -----------------------------------------------------------------------------
SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `accounts`', 'SELECT ''skip accounts''');
-- accounts: App 账号（登录凭证）
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `user_infos`', 'SELECT ''skip user_infos''');
-- user_infos: 用户基础资料/钱包余额等
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `user_exts`', 'SELECT ''skip user_exts''');
-- user_exts: 用户扩展（关注数/靓号等）
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `app_tokens`', 'SELECT ''skip app_tokens''');
-- app_tokens: App 登录 Token
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `cms_tokens`', 'SELECT ''skip cms_tokens''');
-- cms_tokens: CMS 登录 Token（清后需重新登录后台）
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `user_login_devices`', 'SELECT ''skip user_login_devices''');
-- user_login_devices: 用户登录设备
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `user_max_ids`', 'SELECT ''skip user_max_ids''');
-- user_max_ids: 用户号段/最大 ID 记录
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `user_cumulative_stats`', 'SELECT ''skip user_cumulative_stats''');
-- user_cumulative_stats: 用户累计统计
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `user_recharge_cfg_first_recharges`', 'SELECT ''skip user_recharge_cfg_first_recharges''');
-- user_recharge_cfg_first_recharges: 用户档位首充记录
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

-- -----------------------------------------------------------------------------
-- 【删除】货币 / 充值订单
-- -----------------------------------------------------------------------------
SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `currency_logs`', 'SELECT ''skip currency_logs''');
-- currency_logs: 金币/钻石流水
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `recharge_orders`', 'SELECT ''skip recharge_orders''');
-- recharge_orders: 充值订单
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

-- -----------------------------------------------------------------------------
-- 【删除】Dashboard / 运营统计
-- -----------------------------------------------------------------------------
SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `system_total_stats`', 'SELECT ''skip system_total_stats''');
-- system_total_stats: 系统总数据（金币/充值/注册等累计）
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `daily_login_stats`', 'SELECT ''skip daily_login_stats''');
-- daily_login_stats: 日活/日注册/日充值等汇总
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `weekly_login_stats`', 'SELECT ''skip weekly_login_stats''');
-- weekly_login_stats: 周汇总统计
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `monthly_login_stats`', 'SELECT ''skip monthly_login_stats''');
-- monthly_login_stats: 月汇总统计
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `daily_user_logins`', 'SELECT ''skip daily_user_logins''');
-- daily_user_logins: 用户日登录去重明细
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `weekly_user_logins`', 'SELECT ''skip weekly_user_logins''');
-- weekly_user_logins: 用户周登录去重明细
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `monthly_user_logins`', 'SELECT ''skip monthly_user_logins''');
-- monthly_user_logins: 用户月登录去重明细
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `daily_user_recharges`', 'SELECT ''skip daily_user_recharges''');
-- daily_user_recharges: 用户日充值去重明细
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `weekly_user_recharges`', 'SELECT ''skip weekly_user_recharges''');
-- weekly_user_recharges: 用户周充值去重明细
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `monthly_user_recharges`', 'SELECT ''skip monthly_user_recharges''');
-- monthly_user_recharges: 用户月充值去重明细
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `daily_user_gold_consumes`', 'SELECT ''skip daily_user_gold_consumes''');
-- daily_user_gold_consumes: 用户日金币消费去重
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `weekly_user_gold_consumes`', 'SELECT ''skip weekly_user_gold_consumes''');
-- weekly_user_gold_consumes: 用户周金币消费去重
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `monthly_user_gold_consumes`', 'SELECT ''skip monthly_user_gold_consumes''');
-- monthly_user_gold_consumes: 用户月金币消费去重
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `daily_user_diamond_consumes`', 'SELECT ''skip daily_user_diamond_consumes''');
-- daily_user_diamond_consumes: 用户日钻石消费去重
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `weekly_user_diamond_consumes`', 'SELECT ''skip weekly_user_diamond_consumes''');
-- weekly_user_diamond_consumes: 用户周钻石消费去重
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `monthly_user_diamond_consumes`', 'SELECT ''skip monthly_user_diamond_consumes''');
-- monthly_user_diamond_consumes: 用户月钻石消费去重
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `daily_user_audiences`', 'SELECT ''skip daily_user_audiences''');
-- daily_user_audiences: 用户日有效观众去重
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `weekly_user_audiences`', 'SELECT ''skip weekly_user_audiences''');
-- weekly_user_audiences: 用户周有效观众去重
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `monthly_user_audiences`', 'SELECT ''skip monthly_user_audiences''');
-- monthly_user_audiences: 用户月有效观众去重
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `sys_resource_metrics`', 'SELECT ''skip sys_resource_metrics''');
-- sys_resource_metrics: 资源监控细采样
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `sys_resource_metric_aggs`', 'SELECT ''skip sys_resource_metric_aggs''');
-- sys_resource_metric_aggs: 资源监控粗聚合
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

-- -----------------------------------------------------------------------------
-- 【删除】公会 / 直播间 / 场次 / 收益
-- -----------------------------------------------------------------------------
SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `live_guilds`', 'SELECT ''skip live_guilds''');
-- live_guilds: 公会
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `live_guild_transfer_infos`', 'SELECT ''skip live_guild_transfer_infos''');
-- live_guild_transfer_infos: 公会收款/转账资料
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `live_rooms`', 'SELECT ''skip live_rooms''');
-- live_rooms: 直播间实体（通常 ID=主播用户）
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `live_room_cfgs`', 'SELECT ''skip live_room_cfgs''');
-- live_room_cfgs: 每个直播间的房间配置（随房间走，非全局配置）
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `live_room_onlines`', 'SELECT ''skip live_room_onlines''');
-- live_room_onlines: 直播间在线状态
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `live_room_game_recommends`', 'SELECT ''skip live_room_game_recommends''');
-- live_room_game_recommends: 直播间游戏推荐
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `live_room_billing_pays`', 'SELECT ''skip live_room_billing_pays''');
-- live_room_billing_pays: 私密房/计费观看扣费记录
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `live_records`', 'SELECT ''skip live_records''');
-- live_records: 直播场次记录
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `live_record_users`', 'SELECT ''skip live_record_users''');
-- live_record_users: 场次观众明细
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `live_follows`', 'SELECT ''skip live_follows''');
-- live_follows: 关注关系
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `live_revenue_logs`', 'SELECT ''skip live_revenue_logs''');
-- live_revenue_logs: 直播间收益流水（礼物/弹幕/通话等）
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `live_room_income_totals`', 'SELECT ''skip live_room_income_totals''');
-- live_room_income_totals: 直播间生涯累计收益
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `live_room_income_unsettleds`', 'SELECT ''skip live_room_income_unsettleds''');
-- live_room_income_unsettleds: 直播间未结算收益
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `live_room_income_unsettled_archives`', 'SELECT ''skip live_room_income_unsettled_archives''');
-- live_room_income_unsettled_archives: 直播间未结算归档
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `live_room_income_settleds`', 'SELECT ''skip live_room_income_settleds''');
-- live_room_income_settleds: 直播间已结算收益
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `guild_income_totals`', 'SELECT ''skip guild_income_totals''');
-- guild_income_totals: 公会生涯累计收益
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `guild_income_unsettleds`', 'SELECT ''skip guild_income_unsettleds''');
-- guild_income_unsettleds: 公会未结算收益
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `guild_income_unsettled_archives`', 'SELECT ''skip guild_income_unsettled_archives''');
-- guild_income_unsettled_archives: 公会未结算归档
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `guild_income_settleds`', 'SELECT ''skip guild_income_settleds''');
-- guild_income_settleds: 公会已结算收益
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `anchor_income_settlement_logs`', 'SELECT ''skip anchor_income_settlement_logs''');
-- anchor_income_settlement_logs: 主播结算日志
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `guild_income_settlement_logs`', 'SELECT ''skip guild_income_settlement_logs''');
-- guild_income_settlement_logs: 公会结算/打款日志
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `daily_anchor_effective_lives`', 'SELECT ''skip daily_anchor_effective_lives''');
-- daily_anchor_effective_lives: 主播日有效直播统计
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `daily_guild_effective_lives`', 'SELECT ''skip daily_guild_effective_lives''');
-- daily_guild_effective_lives: 公会日有效直播统计
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

-- -----------------------------------------------------------------------------
-- 【删除】通话 / 游戏流水
-- -----------------------------------------------------------------------------
SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `call_orders`', 'SELECT ''skip call_orders''');
-- call_orders: 1v1/直播间通话订单
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `call_users`', 'SELECT ''skip call_users''');
-- call_users: 通话用户状态
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `game_bet_logs`', 'SELECT ''skip game_bet_logs''');
-- game_bet_logs: 游戏下注流水
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `game_win_logs`', 'SELECT ''skip game_win_logs''');
-- game_win_logs: 游戏派奖流水
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

-- -----------------------------------------------------------------------------
-- 【删除】短视频业务
-- -----------------------------------------------------------------------------
SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `short_videos`', 'SELECT ''skip short_videos''');
-- short_videos: 短视频作品
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `short_video_stats`', 'SELECT ''skip short_video_stats''');
-- short_video_stats: 短视频计数（赞/播等）
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `short_video_watches`', 'SELECT ''skip short_video_watches''');
-- short_video_watches: 短视频观看/购买记录
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `short_video_author_stats`', 'SELECT ''skip short_video_author_stats''');
-- short_video_author_stats: 作者维度统计
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `short_video_author_settlement_logs`', 'SELECT ''skip short_video_author_settlement_logs''');
-- short_video_author_settlement_logs: 短视频作者结算日志
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

-- -----------------------------------------------------------------------------
-- 【删除】消息（用户侧会话/未读；保留 activity_messages 配置内容）
-- -----------------------------------------------------------------------------
SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `user_messages`', 'SELECT ''skip user_messages''');
-- user_messages: 私信消息
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `user_message_sessions`', 'SELECT ''skip user_message_sessions''');
-- user_message_sessions: 私信会话
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `user_message_unreads`', 'SELECT ''skip user_message_unreads''');
-- user_message_unreads: 私信未读汇总
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `user_message_unread_details`', 'SELECT ''skip user_message_unread_details''');
-- user_message_unread_details: 私信未读明细
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `user_activity_messages`', 'SELECT ''skip user_activity_messages''');
-- user_activity_messages: 用户活动消息收件箱
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `user_personal_system_messages`', 'SELECT ''skip user_personal_system_messages''');
-- user_personal_system_messages: 个人系统消息
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET @sql := IF(@SAFETY_OK = 1, 'DROP TABLE IF EXISTS `user_system_message_unreads`', 'SELECT ''skip user_system_message_unreads''');
-- user_system_message_unreads: 系统消息未读
PREPARE s FROM @sql; EXECUTE s; DEALLOCATE PREPARE s;

SET FOREIGN_KEY_CHECKS = 1;

SELECT IF(@SAFETY_OK = 1,
  'DONE: business/user/stat tables dropped (config kept). Restart Go to recreate empty tables.',
  'DRY-RUN: no tables dropped. Set @SAFETY_OK := 1 near the top, then re-run.'
) AS result;
