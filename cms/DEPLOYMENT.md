# XR Game Server 管理后台部署文档

## 项目概述

- 项目名称: XR Game Server 管理后台
- 技术栈: Vue 3 + TypeScript + Element Plus + Vue Router + Axios
- 构建工具: Vite

## 环境配置

### 测试环境

- API Base URL: `http://127.0.0.1:8898`
- 构建命令: `npm run build:dev` 或 `npm run build`
- 部署路径: `/cms/`

### 生产环境

- API Base URL: `http://127.0.0.1:1000`
- 构建命令: `npm run build:prod`
- 部署路径: `/cms/`

## 构建说明

### 开发环境构建

```bash
npm run build:dev
```

### 生产环境构建

```bash
npm run build:prod
```

### 通用构建

```bash
npm run build
```

## Nginx 部署配置

将构建后的文件部署到nginx的/cms路径下，请使用以下配置示例：

```nginx
# XR Game Server 管理后台 Nginx 配置示例

server {
    listen 80;
    server_name localhost;  # 修改为你的域名或IP

    # 静态资源 - 前端构建结果
    location /cms/ {
        alias /path/to/your/dist/;  # 修改为实际的dist目录路径
        try_files $uri $uri/ /cms/index.html;  # 支持Vue Router的history模式

        # 带 hash 的 js/css 等可长期缓存；index.html 必须单独禁止缓存
        location ~* \.(js|css|woff2?|ttf|eot|svg|png|jpg|jpeg|gif|webp|ico)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }

        location = /cms/index.html {
            add_header Cache-Control "no-cache, no-store, must-revalidate";
            add_header Pragma "no-cache";
            add_header Expires "0";
        }
        
        # 防止访问隐藏文件（如.git等）
        location ~ /\. {
            deny all;
        }
    }

    # API请求代理到后端服务
    location /api/ {
        proxy_pass http://127.0.0.1:8898/;  # 修改为实际的后端地址
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 其他API请求（根据实际后端API路径调整）
    location /auth/ {
        proxy_pass http://127.0.0.1:8898/;  # 修改为实际的后端地址
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /account/ {
        proxy_pass http://127.0.0.1:8898/;  # 修改为实际的后端地址
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 错误页面
    error_page 500 502 503 504 /50x.html;
    location = /50x.html {
        root html;
    }
}
```

## 部署步骤

1. 构建项目：`npm run build`
2. 将 `dist` 目录中的所有文件上传到服务器的对应目录
3. 配置nginx，将 `/cms/` 路径映射到上传的文件目录
4. 重启nginx服务
5. 访问 `http://your-domain/cms/` 即可使用管理后台