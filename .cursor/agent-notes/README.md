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
| `直播测试服-h5` | h5-live @ 54.241.124.37 | H5/SFTP |

## 2026-08-28 代码丢失（摘要）

- **原因**：清理 `go-src/.../` 误用 `rd /s /q`，路径转义错误，永久删除 `.git` 与大量源码；回收站无备份。
- **恢复**：GitHub reset 到 `0abdb09` + Cursor History + 聊天记录重做 8 个 commit + 未提交改动。
- **教训**：见 `workspace-safety.mdc`；**做完务必 `git push`**。

任务结束后若有新事实，简短更新本文件或新增专题 md。
