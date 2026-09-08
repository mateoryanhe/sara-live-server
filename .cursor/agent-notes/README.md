# Agent 任务上下文（处理任务前先读）

> 本目录 + `.cursor/rules/` 存放运维、约定、事故备忘。

## 规则索引（`.cursor/rules/`）

| 文件 | 内容 |
|------|------|
| [workspace-safety.mdc](../rules/workspace-safety.mdc) | **本机误删、Git push、代码恢复**（2026-08-28 事故） |
| [go-build-output.mdc](../rules/go-build-output.mdc) | Go 编译 `-o` 路径，勿生成 `...` 目录 |
| [server-safety.mdc](../rules/server-safety.mdc) | 远程/数据库删除须问用户 |
| [read-agent-notes.mdc](../rules/read-agent-notes.mdc) | 任务前必读入口 |

## SSH

| Host | 用户 | 用途 |
|------|------|------|
| `直播测试服` | ec2-user @ 54.241.124.37 | Go、MariaDB |
| `直播正式服` | ec2-user @ 52.9.70.64 | 正式 Go |
| `直播审核服` | ec2-user @ 18.144.165.177 | 审核专用（~1GB RAM） |
| `直播测试服-h5` | h5-live @ 54.241.124.37 | H5/SFTP |

## 审核服 MariaDB（2026-08-28）

- **机器**：1 vCPU / ~957MiB RAM / 30G 盘；已加 **512MiB swap**（`vm.swappiness=10`）
- **MariaDB 10.11**，配置 `/etc/my.cnf.d/99-sara-live-review.cnf`
- **端口** `13307`，仅 `127.0.0.1`；库 `live_db`；root 密码与测试/dev 一致
- **内存**：`innodb_buffer_pool_size=96M`，`max_connections=20`，`performance_schema=OFF`，无 binlog
- Go 连接串（与 dev 同格式）：`mysql:root:***@tcp(127.0.0.1:13307)/live_db`

## 审核服磁盘与 /tmp（2026-08-29）

- **`/tmp` 为 tmpfs ~479MB**（约半内存），**勿用于 800MB+ 传输**（会写满失败）
- 大文件/压缩包：用 **`/home/ec2-user/staging`** 中转，解压到目标如 **`/home/ec2-user/cdn/images`**
- 根盘 `/` 约 30G，剩余空间充足；**不宜扩大 /tmp**（受内存限制，扩了易 OOM）
- 已改脚本：`go-review/deploy.bat`、`cms-review/upload.bat`、`avatars-审核服/upload.bat` 均走 staging

## 审核服 cdn/images 被秒清（2026-08-31）

- **现象**：解压/还原 `/home/ec2-user/cdn/images` 后约 1 分钟内变空（非人工 `rm`）。
- **根因**：CMS/日志导出清理与上传资源**共用** `storage_path`，按 mtime 扫目录误删。
- **修复（2026-09-01）**：统一 `module/fileexport`；每个导出文件 `xrtimer.AddOnce(TTL)` 定点删除；前端 `deleteExport` 主动删并取消 timer；**不再扫目录**。重启丢 timer 可接受。
- **解压注意**：包内路径是 `images/...`，应解到 `/home/ec2-user/cdn/` 或用 `staging/restore-review-images.sh`。

## 任务栏一键部署（Win11）

- 入口用 **ASCII** 脚本：`go-test/deploy.bat`、`cms-test/upload.bat`（勿把中文 `一键部署.bat` 编进 launcher，cmd 编码会乱码）
- `taskbar-launcher` 用 `cmd /k`：结束后窗口保持打开，方便看成功/失败；发版 bat **无 pause**（不必回车，关窗即可）
- 已去 pause：`go-test`/`go-prod`/`go-review` 的 `deploy.bat`、`cms-test/upload.bat`，以及对应 `一键部署.bat`
- 固定：`pub-tool/go-test/pin-to-taskbar.bat` 或 `pub-tool/pin-all-taskbar.bat`，再拖桌面快捷方式到任务栏

## 本地 Go / CMS（2026-09-07）

- Go：`pub-tool/go-local/start.bat`（或 `一键启动.bat`）— 停 9443 → `go build` 到 `go-build/xr-game-server.exe` → 以 `config/local` 为 cwd 启动（读 `config/local/config.yaml`）
- CMS：`pub-tool/cms-local/start.bat` — `npm run dev`（5173）

## 2026-08-28 代码丢失（摘要）

- **原因**：清理 `go-src/.../` 误用 `rd /s /q`，路径转义错误，永久删除 `.git` 与大量源码；回收站无备份。
- **恢复**：GitHub reset 到 `0abdb09` + Cursor History + 聊天记录重做 8 个 commit + 未提交改动。
- **教训**：见 `workspace-safety.mdc`；**做完务必 `git push`**。

## module/stat 队列（2026-09-04）

- 无界 `gqueue` + `event.HotStart` 后 `xrtimer` 每秒非阻塞拉空；生产端只 `Push`，校验在消费端
- 登录：`SetLastLoginTime` 同步；日/周/月统计入队
- 观众：`gameevent.ValidAudienceEvent`（`join_room` Pub）
- 热重启：`PrepareRestart` 排空 + `hotrestart.RegisterStatQueueFlush` 与 syndb 同阶段刷盘

## 货币流水 reason 多语言（2026-09-08）

- 文案：`constants/currency/reason_text.go`（zh-CN / zh-TW / en / id；其它语言回落 en）
- App：`POST /userInfo/getCurrencyLog` 返回 `reason` + `reasonText`（按请求头 `Accept-Language`）
- CMS：列表同理；axios 自动带当前界面语言的 `Accept-Language`；异步导出 payload 带 `lang`


- 验证码发信走 CF Email Service REST：`POST /accounts/{account_id}/email/sending/send`（非 AWS SES / 非 CloudFront）
- CMS 表 `cf_email_cfgs`：enabled / accountId / apiToken / fromEmail（如 `noreply@mail.saralive.net`）
- 入口：`POST /auth/sendEmailCode`、`POST /auth/emailLogin`（免鉴权）、`POST /auth/bindEmail`（需登录）
- 渠道 `EmailChannel=7`；绑定写 `user_exts.email`
- 占用规则：仅看未注销账号（EmailChannel open_id 或已绑定 email）；注销即释放；UnCancel 再检查邮箱
- `user_exts.email` → userId 走 `emailUserIdCacheMgr`（绑定 Publish / 注销 Invalidate）
- 验证码：发信成功写入 `emailVerifyCodeCache`（5 分钟）；登录/绑定校验后删除；调试码 `981200`
- 发信限流（双 gcache，不入库）：`emailSendCooldownCache` 1 分钟；`emailSendDailyCache` 每日 10 次、TTL 到本地 0 点
- 菜单：配置 → Cloudflare邮件（`/cfEmail`）

## 工会菜单分组（2026-09-07）

- 运营 → 工会管理 拆三层：`基础管理` / `账号权限` / `数据流水`
- 权限树 `PERMISSION_MENU_TREE` 同步嵌套；页面 module 名不变（已授权角色无需改库）
- `GuildManagement` 菜单文案改为「工会列表」（避免与父级「工会管理」重复）

- 渠道 `CoinMerchantChannel=8`；`UserTypeCoinMerchant=6`（不参与系统统计）
- 账号仍走 `accounts` + `user_infos` 缓冲；用户名=`open_id`，密码 MD5 存 `accounts.password`
- CMS：用户管理 → 币商；接口 `/coinMerchant/*`（列表/新建/重置密码/注销）
- App 登录：`POST /auth/coinMerchantLogin`（免鉴权；username+password；不自动注册；token=`userId.token`）
- 充值档位：表 `coin_merchant_recharge_cfgs`（name/USD price/gold/status）；写库后整体刷 `atomic` 上架缓存；CMS 在「充值会员 → 币商充值档位」
- App 查档位：`POST /coinMerchantRechargeCfg/coinMerchantRechargeCfgListForApp`（需登录，仅上架缓存）
- App yhpay 下单：`POST /rechargeOrder/createCoinMerchantChannelRechargeOrder`（需登录+`UserTypeCoinMerchant`；`cfgId`+`currencyCode=IDR`；`payChannel=4`；无首充加赠；回调复用 `/webhook/yhpay/payin`）
- 轮询成功：`POST /rechargeOrder/checkRechargeOrderSuccess`（与普通充值相同）
- App 转赠金币：`POST /gold/transferGold`（仅币商；`targetUserId`+`amount` 最多2位小数；扣币商加目标用户；流水 reason 32转出/33收入）

- 优先读请求头 `CF-IPCountry`（如 `US`→`美国`），无效/`XX`/`T1` 再回退 GeoLite；入库存中文名

## Flutter / Android 本机工具链（2026-09-07）

- 根目录：`D:\tools`
  - Flutter 3.47.2：`D:\tools\flutter`
  - JDK 17（Temurin，打 APK 用，勿用本机 JDK 25）：`D:\tools\jdk-17`
  - Android SDK：`D:\tools\android-sdk`（platform 34/35/36、build-tools、NDK 28.2）
- 中国镜像：`PUB_HOSTED_URL=https://pub.flutter-io.cn`，`FLUTTER_STORAGE_BASE_URL=https://storage.flutter-io.cn`
- 会话加载：`. D:\tools\env-flutter.ps1` 或 `D:\tools\env-flutter.bat`（会把 `127.0.0.1:port` 代理补成 `http://`）
- 重装脚本：`D:\tools\install-flutter-android.ps1`；补包：`D:\tools\install-android-packages.ps1`
- App 工程：`flutter-client/`（已 gitignore）；打 APK：`flutter-client\build-apk.bat`
- 注意：本机 `HTTP_PROXY` 若无协议，sdkmanager 会挂；VS 未装不影响 Android APK

## CMS 本地启动（2026-09-07）

- 脚本：`pub-tool/cms-local/start.bat`（或 `一键启动.bat`）
- 目录：`cms/`，`npm run dev`，默认 http://127.0.0.1:5173
- API：`.env.development` → `VITE_API_BASE_URL=http://127.0.0.1:9443`（需先起本地 Go）

## 直播间 1v1 通话充值门槛（2026-09-07）

- `liveRoomCall`：**暂不**校验累计充值满 10 USD（`CanInitiateLiveRoomCall`）
- `AllowCallIcon` / `ViewerCanUseLiveRoomCall`：仍校验 10 USD（客户端入口隐藏）
- 常量：`LiveRoomCallMinTotalRechargeUSD = 10`

## 工会可见性（2026-09-07）

- 表：`live_guild_visibilities`（`guild_id` + `cms_user_id`，唯一索引）
- 创建工会：自动写入当前登录 CMS 用户一条可见记录
- 列表 `guildList` / 下架列表：非管理员只看可见表；管理员/超管看全部
- CMS 页：运营 → 工会 → **工会可见性**（侧栏靠**角色授权**显示，不写死管理员）
  - 选 CMS 用户 → `el-transfer` 双栏：左=未授权上架工会，右=已授权；箭头移动即批量授权/撤销
  - 接口：`guildVisibilityByUserList`、`batchGrantGuildVisibility`、`batchRevokeGuildVisibility`、`guildListForVisibility`（拉全量上架，不过滤可见性表）
  - 使用场景上通常只把该页授给管理员，但机制与其它 CMS 页一致
- **社交日志同规则**（非管理员仅可见工会下主播数据；管理员看全部）：
  - 直播记录 / 收益流水 / 每日流水 / 本周流水 / 视频通话日志（含对应导出）
  - 关联：`live_rooms.guild_id`（通话按 `receiver_id`=主播间）
  - 辅助包：`module/cmsvis`（避免与 liveroom/liverecord 循环依赖）
  - **未限制**：短视频观看记录（按观众维度，无自然工会字段）
- 历史工会无记录时普通用户看不到（可用本页授权，或手工插表）

