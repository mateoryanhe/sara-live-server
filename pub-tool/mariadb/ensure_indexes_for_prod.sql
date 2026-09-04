-- =============================================================================
-- live_db 上线前必要索引补齐（与 go-src/entity gorm 标签对齐，2026-09-04）
-- =============================================================================
-- 用途：
--   1) 清空业务表后重启 Go，AutoMigrate 会按 entity 自动建这些索引；本脚本可作核对清单
--   2) 若表已存在且不想靠 AutoMigrate 补索引，可手动执行本脚本
--
-- 说明：
--   - MariaDB 10.5+ 支持 CREATE INDEX IF NOT EXISTS
--   - third_order_id / order_id 等可能有大量空串，用普通 INDEX，不用 UNIQUE
--   - AutoMigrate 不会 DROP 旧索引；call_orders 等噪音索引需本脚本末尾按需 DROP
-- =============================================================================

SELECT DATABASE() AS current_database;

-- -----------------------------------------------------------------------------
-- MUST：登录 / 流水 / 支付 / 消息
-- -----------------------------------------------------------------------------

-- accounts: 登录 channel + open_id
CREATE INDEX IF NOT EXISTS `idx_account_channel_openid` ON `accounts` (`channel`, `open_id`);

-- currency_logs: 富豪榜(type,action,时间) + 用户流水(user,type,时间)
CREATE INDEX IF NOT EXISTS `idx_cl_type_action_created` ON `currency_logs` (`type`, `action`, `created_at`);
CREATE INDEX IF NOT EXISTS `idx_cl_user_type_created` ON `currency_logs` (`user_id`, `type`, `created_at`);

-- live_revenue_logs: 主播榜 / 贡献榜 / 通话退款查找
CREATE INDEX IF NOT EXISTS `idx_lrl_status_created` ON `live_revenue_logs` (`status`, `created_at`);
CREATE INDEX IF NOT EXISTS `idx_lrl_room_status_created` ON `live_revenue_logs` (`room_id`, `status`, `created_at`);
CREATE INDEX IF NOT EXISTS `idx_lrl_biz_sender_type_status` ON `live_revenue_logs` (`biz_id`, `sender_id`, `revenue_type`, `status`);

-- recharge_orders: 支付回调幂等 + 用户订单列表
CREATE INDEX IF NOT EXISTS `idx_recharge_orders_third_order_id` ON `recharge_orders` (`third_order_id`);
CREATE INDEX IF NOT EXISTS `idx_ro_user_status` ON `recharge_orders` (`user_id`, `status`);

-- user_message_unread_details: 会话列表 user + mutual + updated_at
CREATE INDEX IF NOT EXISTS `idx_umud_user_mutual_updated` ON `user_message_unread_details` (`user_id`, `mutual_chat`, `updated_at`);

-- user_message_sessions: 私信历史 session + is_deleted
CREATE INDEX IF NOT EXISTS `idx_ums_session_deleted` ON `user_message_sessions` (`session_id`, `is_deleted`);

-- -----------------------------------------------------------------------------
-- SHOULD：用户 / 游戏 / 直播 / 通话 / 结算 / 短视频
-- -----------------------------------------------------------------------------

CREATE INDEX IF NOT EXISTS `idx_user_infos_last_login_time` ON `user_infos` (`last_login_time`);
CREATE INDEX IF NOT EXISTS `idx_user_exts_cancel_code` ON `user_exts` (`cancel_code`);

CREATE INDEX IF NOT EXISTS `idx_game_bet_created` ON `game_bet_logs` (`created_at`);

CREATE INDEX IF NOT EXISTS `idx_live_records_start_time` ON `live_records` (`start_time`);
CREATE INDEX IF NOT EXISTS `idx_live_room_status_updated` ON `live_rooms` (`status`, `updated_at`);

CREATE INDEX IF NOT EXISTS `idx_call_caller_start` ON `call_orders` (`caller_id`, `call_start_time`);
CREATE INDEX IF NOT EXISTS `idx_call_receiver_start` ON `call_orders` (`receiver_id`, `call_start_time`);

CREATE INDEX IF NOT EXISTS `idx_gis_guild_status_created` ON `guild_income_settlement_logs` (`guild_id`, `status`, `created_at`);
CREATE INDEX IF NOT EXISTS `idx_ais_room_created` ON `anchor_income_settlement_logs` (`room_id`, `created_at`);

CREATE INDEX IF NOT EXISTS `idx_short_video_watches_video_id` ON `short_video_watches` (`video_id`);
CREATE INDEX IF NOT EXISTS `idx_short_video_watches_paid_time` ON `short_video_watches` (`paid_time`);

-- -----------------------------------------------------------------------------
-- OPTIONAL：清理噪音索引（仅在确认无业务依赖后执行；默认注释）
-- GORM AutoMigrate 不会自动 DROP，清空重建表则无需手工删
-- -----------------------------------------------------------------------------
-- ALTER TABLE `call_orders` DROP INDEX `idx_call_orders_answer_time`;
-- ALTER TABLE `call_orders` DROP INDEX `idx_call_orders_caller_confirm_time`;
-- ALTER TABLE `call_orders` DROP INDEX `idx_call_orders_receiver_confirm_time`;
-- ALTER TABLE `call_orders` DROP INDEX `idx_call_orders_caller_heart_time`;
-- ALTER TABLE `call_orders` DROP INDEX `idx_call_orders_receiver_heart_time`;
-- ALTER TABLE `call_orders` DROP INDEX `idx_call_orders_caller_hang_up_time`;
-- ALTER TABLE `call_orders` DROP INDEX `idx_call_orders_receiver_hang_up_time`;
-- ALTER TABLE `call_orders` DROP INDEX `idx_call_orders_order_end_time`;
-- ALTER TABLE `call_orders` DROP INDEX `idx_call_orders_charge_time`;
-- ALTER TABLE `call_users` DROP INDEX `idx_call_users_heart_time`;
-- ALTER TABLE `live_room_billing_pays` DROP INDEX `idx_live_room_billing_pays_free_time`;
-- ALTER TABLE `user_exts` DROP INDEX `idx_user_exts_pretty_id`;
-- ALTER TABLE `user_exts` DROP INDEX `idx_user_exts_cancel_code_expire_at`;
-- ALTER TABLE `recharge_orders` DROP INDEX `idx_recharge_orders_cfg_id`;
