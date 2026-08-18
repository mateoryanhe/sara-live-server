-- =============================================================================
-- 内测清空：用户 + 主播 + 工会业务数据（DROP TABLE）
-- =============================================================================
-- 用途：清空 App 用户、主播、直播间、收益、流水、统计等业务数据，便于内测重来。
--
-- ⚠️  警告：
--   1. 不可恢复，仅在测试库执行！
--   2. 不删除 CMS 账号、角色权限、礼物/配置类表。
--   3. 执行后请重启 xr-game-server，AutoMigrate 会自动重建空表。
--   4. 建议执行前备份：mysqldump -u... -p... 库名 > backup.sql
-- =============================================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- -----------------------------------------------------------------------------
-- 1. 用户消息 / 活动消息（用户侧）
-- -----------------------------------------------------------------------------
DROP TABLE IF EXISTS `user_message_unread_details`;
DROP TABLE IF EXISTS `user_message_unreads`;
DROP TABLE IF EXISTS `user_message_sessions`;
DROP TABLE IF EXISTS `user_messages`;
DROP TABLE IF EXISTS `user_system_message_unreads`;
DROP TABLE IF EXISTS `user_personal_system_messages`;
DROP TABLE IF EXISTS `user_activity_messages`;

-- -----------------------------------------------------------------------------
-- 2. 通话 / 游戏 / 充值订单（用户行为流水）
-- -----------------------------------------------------------------------------
DROP TABLE IF EXISTS `call_orders`;
DROP TABLE IF EXISTS `call_users`;
DROP TABLE IF EXISTS `game_bet_logs`;
DROP TABLE IF EXISTS `game_win_logs`;
DROP TABLE IF EXISTS `recharge_orders`;
DROP TABLE IF EXISTS `currency_logs`;

-- -----------------------------------------------------------------------------
-- 3. 短视频（用户产生的内容）
-- -----------------------------------------------------------------------------
DROP TABLE IF EXISTS `short_video_watches`;
DROP TABLE IF EXISTS `short_video_stats`;
DROP TABLE IF EXISTS `short_videos`;

-- -----------------------------------------------------------------------------
-- 4. 直播流水 / 关注 / 在线 / 计费
-- -----------------------------------------------------------------------------
DROP TABLE IF EXISTS `live_revenue_logs`;
DROP TABLE IF EXISTS `live_record_users`;
DROP TABLE IF EXISTS `live_records`;
DROP TABLE IF EXISTS `live_follows`;
DROP TABLE IF EXISTS `live_room_onlines`;
DROP TABLE IF EXISTS `live_room_billing_pays`;
DROP TABLE IF EXISTS `live_private_room_billings`;
DROP TABLE IF EXISTS `live_tickets`;
DROP TABLE IF EXISTS `live_room_game_recommends`;

-- -----------------------------------------------------------------------------
-- 5. 主播 / 工会收益 & 结算
-- -----------------------------------------------------------------------------
DROP TABLE IF EXISTS `anchor_income_settlement_logs`;
DROP TABLE IF EXISTS `guild_income_settlement_logs`;
DROP TABLE IF EXISTS `daily_anchor_effective_lives`;
DROP TABLE IF EXISTS `daily_guild_effective_lives`;

DROP TABLE IF EXISTS `live_room_income_unsettled_archives`;
DROP TABLE IF EXISTS `guild_income_unsettled_archives`;

DROP TABLE IF EXISTS `live_room_income_unsettleds`;
DROP TABLE IF EXISTS `live_room_income_settleds`;
DROP TABLE IF EXISTS `live_room_income_totals`;

DROP TABLE IF EXISTS `guild_income_unsettleds`;
DROP TABLE IF EXISTS `guild_income_settleds`;
DROP TABLE IF EXISTS `guild_income_totals`;

-- -----------------------------------------------------------------------------
-- 6. 直播间 / 工会（主播主体）
-- -----------------------------------------------------------------------------
DROP TABLE IF EXISTS `live_room_cfgs`;
DROP TABLE IF EXISTS `live_rooms`;
DROP TABLE IF EXISTS `live_guilds`;

-- -----------------------------------------------------------------------------
-- 7. 用户账号主体
-- -----------------------------------------------------------------------------
DROP TABLE IF EXISTS `user_cumulative_stats`;
DROP TABLE IF EXISTS `user_login_devices`;
DROP TABLE IF EXISTS `user_exts`;
DROP TABLE IF EXISTS `user_infos`;
DROP TABLE IF EXISTS `app_tokens`;
DROP TABLE IF EXISTS `accounts`;
DROP TABLE IF EXISTS `user_max_ids`;

-- -----------------------------------------------------------------------------
-- 8. 用户相关统计（日 / 周 / 月）
-- -----------------------------------------------------------------------------
DROP TABLE IF EXISTS `daily_user_audiences`;
DROP TABLE IF EXISTS `daily_user_diamond_consumes`;
DROP TABLE IF EXISTS `daily_user_gold_consumes`;
DROP TABLE IF EXISTS `daily_user_recharges`;
DROP TABLE IF EXISTS `daily_user_logins`;
DROP TABLE IF EXISTS `daily_login_stats`;

DROP TABLE IF EXISTS `weekly_user_audiences`;
DROP TABLE IF EXISTS `weekly_user_diamond_consumes`;
DROP TABLE IF EXISTS `weekly_user_gold_consumes`;
DROP TABLE IF EXISTS `weekly_user_recharges`;
DROP TABLE IF EXISTS `weekly_user_logins`;
DROP TABLE IF EXISTS `weekly_login_stats`;

DROP TABLE IF EXISTS `monthly_user_audiences`;
DROP TABLE IF EXISTS `monthly_user_diamond_consumes`;
DROP TABLE IF EXISTS `monthly_user_gold_consumes`;
DROP TABLE IF EXISTS `monthly_user_recharges`;
DROP TABLE IF EXISTS `monthly_user_logins`;
DROP TABLE IF EXISTS `monthly_login_stats`;

DROP TABLE IF EXISTS `system_total_stats`;

SET FOREIGN_KEY_CHECKS = 1;

-- =============================================================================
-- 保留不动（配置 / CMS，无需重建）：
--   cms_users, cms_roles, permissions, cms_tokens
--   live_gifts, live_cfgs, live_room_tags, anchor_salary_cfgs, live_revenue_share_cfgs
--   home_banners, vip_cfgs, agora_cfgs, aliyun_text_moderation_cfgs
--   recharge_cfgs, google_play_cfgs, game_cfgs, game_platform_cfgs
--   account_cfgs, random_nicknames, wallet_exchange_cfgs, customer_service_cfgs
--   privacy_policy_cfgs, short_video_cfgs, short_video_categories
--   activity_messages（活动消息模板）
--   upload_resource_cfgs, preload_cfgs, data_sync_cfgs, simulator_cpu_keywords, app_pkgs
-- =============================================================================
--
-- 执行完成后：
--   1. 重启 go 服务（AutoMigrate 重建上表）
--   2. 清空 Redis 缓存（如有用户/直播间/礼物缓存）
-- =============================================================================
