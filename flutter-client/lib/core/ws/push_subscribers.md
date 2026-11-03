# WS 推送订阅（AppBus）

- **连接**：`WsClient.onState` → [AppScope] 只发布 [WsConnectionBusEvent]（`connecting` / `connected` / `disconnected`）。
- **下行 push**：`WsClient` 帧解析 → [WsBusEvent]（服务端 `cmd` + `data`）。
- **各模块自行 `appBus.add`**，互不调用。约定见 `.cursor/rules/flutter-push-subscribe.mdc`。

调试：`HTTP_LOG_ENABLED=true` 时 [WsLogger] 与 [HttpLogger] 共用 sink。

| 层级 | 职责 |
|------|------|
| **Repository** | `invalidate` / patch / **WS 上线 HTTP 补偿**；不碰 UI |
| **ViewModel / Controller** | 同一事件再订一份 → 读 Repository 或 HTTP → 刷新本模块 UI |

## 连接事件 [WsConnectionBusEvent]

| state | 含义 | 典型订阅方 |
|-------|------|------------|
| `connecting` | 正在连 `/ws` | （可选，一般忽略） |
| `connected` | 已上线，可收 push | `UserRepository.refreshWalletFromServer`；各 VM 刷 UI |
| `disconnected` | 离线 / 断线 | （可选） |
| `resumed`（App 生命周期） | 回前台 | 同上：后台偷偷 `/userInfo/get` 更新 profile 缓存 |

## 下行 cmd（[WsBusEvent]）

| cmd | 含义 | Repository | ViewModel / Controller |
|-----|------|------------|------------------------|
| 8 | 直播间送礼 | — | `LiveRoomViewerController`, `LiveRoomAnchorController` |
| 9 | 公屏文字 | — | 同上 |
| 22 | 付费弹幕 | — | 同上 |
| 10 | 钻石余额 | `UserRepository` patch profile | `MineViewModel`；`LiveRoomViewerController` |
| 11 | 金币余额 | 同上 | 同上 |
| 40 | 充值成功 | `UserRepository` patch `goldBalance` | 同上 |
| 12 | VIP 等级 | `UserRepository` / `VipRepository` invalidate | `MineViewModel`；`VipViewModel` |
| 23 | 主播下播 | — | `LiveRoomViewerController` |
| 24 | 观众列表 | — | `LiveRoomAnchorController`；`AnchorAudienceViewModel` |
| 25–32 | 1v1 通话 | `OneToOneRepository` 等 | `CallController`；`OneToOneViewModel` |

`connected` 时列表/主播缓存补偿：`LiveRoomRepository`、`OneToOneRepository`、`VipRepository` 等 **各自 invalidate**，对应 VM **各自 refresh**。

新增能力：**不要**在 `AppScope` 里集中 switch；在 Repository + 需要展示的 VM/Controller 各加订阅。
