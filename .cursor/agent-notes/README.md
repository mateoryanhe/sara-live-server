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

## 2026-08-28 代码丢失（摘要）

- **原因**：清理 `go-src/.../` 误用 `rd /s /q`，路径转义错误，永久删除 `.git` 与大量源码；回收站无备份。
- **恢复**：GitHub reset 到 `0abdb09` + Cursor History + 聊天记录重做 8 个 commit + 未提交改动。
- **教训**：见 `workspace-safety.mdc`；**做完务必 `git push`**。

任务结束后若有新事实，简短更新本文件或新增专题 md。
