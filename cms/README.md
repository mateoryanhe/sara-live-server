# xr-web

## 项目介绍

基于 Vue3 + TypeScript + Element Plus + Vite 构建的管理后台前端项目。

## 技术栈

- Vue 3
- TypeScript
- Element Plus
- Vue Router
- Axios
- Vite

## 目录结构

```
src/
├── api/                    # API请求模块
│   ├── modules/            # 各业务模块API
│   │   ├── auth.ts         # 权限相关API
│   │   ├── account.ts      # 账号相关API
│   │   ├── agent.ts        # 智能体相关API
│   │   ├── coinCfg.ts      # 币种配置相关API
│   │   └── globalCfg.ts    # 全局配置相关API
│   └── request.ts          # axios请求封装
├── assets/                 # 静态资源
├── components/             # 公共组件
├── config/                 # 配置文件
│   └── env.ts              # 环境配置
├── router/                 # 路由配置
├── types/                  # TypeScript类型定义
│   └── api.ts              # API类型定义
├── views/                  # 页面组件
│   ├── account/            # 账号管理页面
│   ├── agent/              # 智能体管理页面
│   ├── coin/               # 币种管理页面
│   ├── config/             # 系统配置页面
│   ├── dashboard/          # 仪表盘页面
│   ├── layout/             # 布局组件
│   └── login/              # 登录页面
├── App.vue                 # 根组件
└── main.ts                 # 入口文件
```

## 环境配置

### 开发环境

- Node.js: >= 16.0.0
- npm: >= 8.0.0

### 环境变量

项目使用 `.env` 文件进行环境配置：

- `.env.development`: 开发环境配置
- `.env.production`: 生产环境配置
- `.env.test`: 测试环境配置

## 安装依赖

```bash
npm install
```

## 开发运行

```bash
# 启动开发服务器
npm run dev

# 构建开发环境版本
npm run build:dev

# 构建生产环境版本
npm run build:prod

# 构建通用版本
npm run build

# 预览构建结果
npm run preview
```

## 部署说明

项目构建后的文件可部署在nginx服务器的 `/cms` 路径下。具体配置请参考 [DEPLOYMENT.md](./DEPLOYMENT.md)。

## API 接口

项目通过 API 模块与后端服务交互，所有请求都经过统一的请求封装处理。

- 权限相关: `/auth/*`
- 账号管理: `/account/*`
- 全局配置: `/globalCfg/*`
- 智能体管理: `/agent/*`
- 币种配置: `/coinCfg/*`

## 路由配置

项目使用 Vue Router 进行路由管理，支持嵌套路由和路由守卫。

## 代码规范

- TypeScript 严格模式
- ESLint + Prettier 代码格式化
- 组件命名采用 PascalCase
- 文件命名采用 kebab-case
- API 接口统一管理

## 注意事项

- 所有 API 请求都需携带认证信息（token、authId）
- 组件按功能模块组织在 views 目录下
- 类型定义统一管理在 types 目录下
- 静态资源统一管理在 assets 目录下

## 部署配置

项目构建后部署在nginx的 `/cms` 目录下，支持history模式的路由。请确保nginx配置正确，将 `/cms/` 路径映射到构建输出的 `dist`
目录。