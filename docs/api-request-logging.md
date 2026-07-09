# API 请求分段日志说明

本文档说明 `go-src/core/httpserver` 中 API 请求的分段日志设计，帮助排查慢接口、鉴权问题、大 body/大响应等问题。

实现代码主要在：

- `go-src/core/httpserver/log_util.go`
- `go-src/core/httpserver/middleware.go`
- `go-src/core/httpserver/app_auth_middleware.go` / `cms_auth_middleware.go`
- `go-src/core/httpserver/api_res_middleware.go` / `cms_res_custom_middleware.go`
- `go-src/core/httpserver/response_middleware_util.go`
- `go-src/core/httpserver/http_server.go`（`HookAfterOutput`）

**耗时单位**：所有分段耗时字段统一为 **毫秒（ms）**，字段名以 `Ms` 结尾或带 `Ms=` 后缀。

---

## 1. 日志链路概览

一次完整的 API 请求，按时间顺序最多输出 **6 条** Info 日志：

```
收到前端请求
    ↓
鉴权完成          （需鉴权路由才有）
    ↓
读取请求Body
    ↓
Handler执行完成
    ↓
应答写入缓冲区
    ↓
应答输出完成      （gzip + 真正写出，HookAfterOutput）
```

### 1.1 请求在中间件中的实际顺序

```
全局: CORS → [Gzip]
    ↓
middlewareLogReq          → 收到前端请求
    ↓
MiddlewareAppAuth / MiddlewareCmsAuth   → 鉴权完成（可选）
    ↓
apiResponseMiddleware / customResponseMiddleware
    ├─ 读取请求Body
    ├─ Handler（Parse + 业务）
    ├─ Handler执行完成
    └─ 应答写入缓冲区
    ↓
[Gzip 压缩] → TCP 写出
    ↓
HookAfterOutput         → 应答输出完成
```

### 1.2 关联同一次请求

从 **「鉴权完成」** 起，后续日志均包含：

| 字段 | 来源 | 说明 |
|------|------|------|
| `reqId` | Header `reqId` | 前端/request 侧请求 ID，**串联日志的主键** |
| `authId` | App: `Authorization` 中 `.` 前一段；CMS: Header `authId` | 用户/操作员 ID |

建议排查时用 `grep reqId=xxx` 拉取同一次请求的全部分段。

> 说明：「收到前端请求」目前不含 `reqId`，高并发下可用 `url + ip + time` 与后续日志人工对应；如需更强关联可后续在该步补充 `reqId`。

---

## 2. 各阶段日志详解

### 2.1 收到前端请求

**触发**：`middlewareLogReq` 入口，**只读 Header，不读 Body**。

**示例**：

```
收到前端请求,time=2026-06-05 17:00:00.150,enterTime=2026-06-05 17:00:00.120,从队列进入到中间件时间间隔Ms=30ms,method=POST,url=/api/liveRoom/join,ip=113.109.204.237,headers={"Authorization":["123456.token"],"reqId":["789"],...}
```

| 字段 | 含义 |
|------|------|
| `time` | 打日志时刻 |
| `enterTime` | GoFrame `r.EnterTime`，请求进入框架的时刻 |
| `从队列进入到中间件时间间隔Ms` | `time - enterTime`（ms），反映排队/全局中间件开销 |
| `method` / `url` / `ip` | HTTP 方法、URI、客户端 IP |
| `headers` | 全部请求头（JSON 字符串） |

---

### 2.2 鉴权完成

**触发**：`MiddlewareAppAuth` / `MiddlewareCmsAuth`，成功或失败都会记录。

**示例**：

```
鉴权完成,time=2026-06-05 17:00:00.152,reqId=789,authId=123456,authMs=2ms
```

| 字段 | 含义 |
|------|------|
| `authMs` | 鉴权逻辑耗时（ms），不含后续业务 |

**无此日志的情况**：`RegNonAuthAPI`、`RegCMSNonAuthCustomizeRes` 等无需鉴权路由。

---

### 2.3 读取请求Body

**触发**：响应中间件在 `Next()`（Handler）**之前**，主动读取/解析 Body 用于日志。

**示例（JSON）**：

```
读取请求Body,time=...,reqId=789,authId=123456,bodyMs=12ms,bodyLength=128,bodyContent={"roomId":1001}
```

**示例（multipart 表单 + 文件混传）**：

```
读取请求Body,...,bodyContent={"title":"视频","tags":"娱乐","video":{"filename":"a.mp4","size":5242880,"skipped":"[文件已省略,不输出内容]"}}
```

| 字段 | 含义 |
|------|------|
| `bodyMs` | 读取/解析 Body 耗时（ms） |
| `bodyLength` | Body 长度；multipart 优先用 `Content-Length` |
| `bodyContent` | 见下方 [Body 记录规则](#3-body-记录规则) |

---

### 2.4 Handler执行完成

**触发**：`r.Middleware.Next()` 返回后，即 Parse + 业务 Handler 执行完毕。

**示例**：

```
Handler执行完成,time=...,reqId=789,authId=123456,handlerMs=85ms,url=/api/liveRoom/join
```

| 字段 | 含义 |
|------|------|
| `handlerMs` | Handler 阶段耗时（ms），含 DTO Parse/校验 + 业务逻辑 |

> 若此阶段慢，需进入具体 Handler 内查 DB、Redis、外部 HTTP 等；HTTP 层日志无法细分到每一行代码。

---

### 2.5 应答写入缓冲区

**触发**：`writeResponse` 内 JSON 序列化并 `r.Response.Write()` 写入 GoFrame 响应缓冲区。

**示例**：

```
应答写入缓冲区,time=...,reqId=789,authId=123456,writeMs=3ms,respBytes=256,url=/api/liveRoom/join
```

| 字段 | 含义 |
|------|------|
| `writeMs` | 序列化 + 写缓冲区耗时（ms） |
| `respBytes` | 写入缓冲区的响应字节数 |

鉴权失败等走 `WriteFailJson` 的路径也会写入缓冲区并打时间戳，供下一步 Hook 使用，但不一定有 Handler/Body 等前置日志。

---

### 2.6 应答输出完成

**触发**：`HookAfterOutput`，响应已从缓冲区经 gzip（若启用）并写出到客户端之后。

**示例**：

```
应答输出完成,time=...,reqId=789,authId=123456,afterOutputMs=8ms,gzip=true,totalMs=528ms,url=/api/liveRoom/join
```

| 字段 | 含义 |
|------|------|
| `afterOutputMs` | 从「应答写入缓冲区」到输出完成的耗时（ms），**含 gzip 压缩 + TCP 写出** |
| `gzip` | 最终响应是否带 `Content-Encoding: gzip` |
| `totalMs` | 从 `EnterTime` 到 `LeaveTime` 的总耗时（ms），**适合脚本筛慢接口** |

---

## 3. Body 记录规则

设计原则：**小体积结构化数据记内容，大文件只记元信息**。

| Content-Type / 场景 | bodyContent | bodyLength |
|---------------------|-------------|------------|
| `application/json` | 完整 JSON 原文 | 字节长度 |
| `application/x-www-form-urlencoded` | 表单字段 JSON | 序列化后长度 |
| `multipart/*`（表单+文件混传） | 文本字段 + 文件 `{filename,size,skipped}` | 优先 `Content-Length` |
| 纯文件（`image/*`、`application/octet-stream`） | `[文件已省略,不输出内容]` | `Content-Length` |
| 非 JSON 且无表单 | `[非JSON报文,已省略内容]` | 原始 body 长度 |
| GET / 无 body | 空 | `0` |

multipart 文本字段会从 `GetFormMap()` 与 `MultipartForm.Value` 双通道采集，避免「有文件时表单字段漏记」。

---

## 4. 哪些路由会打这套日志

| 注册方式 | 鉴权 | Body/Handler/写缓冲/输出 |
|----------|------|---------------------------|
| `RegAPI` | App 鉴权 | 完整链路 |
| `RegNonAuthAPI` | 无 | 无「鉴权完成」，其余有 |
| `RegAppCustomizeRes` | App 鉴权 | 完整链路 |
| `RegCMS` | CMS 鉴权 | 完整链路 |
| `RegCMSCustomizeRes` | CMS 鉴权 | 完整链路 |
| `RegCMSNonAuthCustomizeRes` | 无 | 无「鉴权完成」，其余有 |

**不会打 API 分段日志**：静态文件、WebSocket 等未走 `writeResponse` / `WriteFailJson` 的请求（「应答输出完成」也不会出现）。

---

## 5. 慢接口排查指南

推荐先看 **「应答输出完成」** 中的 `totalMs`，再用同一 `reqId` 回看各段：

| 现象 | 可能原因 | 下一步 |
|------|----------|--------|
| `从队列进入到中间件时间间隔Ms` 大 | 连接排队、服务器负载高 | 看 QPS、goroutine、全局中间件 |
| `authMs` 大 | Token 校验慢 | 查 `xrtoken`、DB/缓存 |
| `bodyMs` 大且 `bodyLength` 大 | 大 body / multipart 解析 | 看上传体积、网络 |
| `handlerMs` 大 | 业务或 DB 慢 | 进 Handler 加子日志 |
| `writeMs` 大且 `respBytes` 大 | 响应体过大 | 优化返回字段/分页 |
| `afterOutputMs` 大 | gzip 或网络写出慢 | 看 `gzip` 标志、响应大小、客户端网络 |

各段耗时关系（近似）：

```
totalMs ≈ 排队 + 鉴权 + 读Body + Handler + 写缓冲 + afterOutput
```

---

## 6. 与现有脚本

仓库 `docs/` 下有过滤脚本，例如：

- `filter_slow_access_requests.sh` — GoFrame **access** 日志（整体请求耗时，单位 ms）
- `filter_api_requests.sh` — 应用 **业务日志**（可按新字段调整）

针对当前分段日志，筛慢接口建议：

```bash
# 按 totalMs 筛「应答输出完成」（示例阈值 400ms）
grep '应答输出完成' app.log | awk -F'totalMs=' '{split($2,a,"ms"); if (a[1]+0 >= 400) print}'

# 按 reqId 拉全链路
grep 'reqId=789' app.log
```

可按团队习惯更新 `filter_api_requests.sh`，匹配 `totalMs`、`handlerMs` 等字段。

---

## 7. 设计要点（给维护者）

1. **Header 阶段不读 Body**，避免无意义的 IO。
2. **Body 日志在 Handler 前读取**，GoFrame 会缓存 `bodyContent`，后续 `Parse` 不会重复读网络。
3. **文件不上屏**，只记 filename/size，控制日志体积。
4. **`afterOutputMs` 补齐 gzip/写出**，避免「写缓冲很快但客户端仍慢」的盲区。
5. **鉴权失败**走 `WriteFailJson` 仍会 stash 写缓冲时间戳，Hook 可记录输出耗时。
6. **耗时统一 ms**，便于与 access 日志、脚本阈值对齐。

---

## 8. 日志级别

当前分段日志均为 **`Info`**。慢接口告警/filter 由外部脚本或日志平台按 `totalMs` 等字段处理，不在代码内硬编码阈值。
