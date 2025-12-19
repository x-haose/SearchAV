# SearchAV

视频搜索聚合服务，Go 后端 + SvelteKit 前端。

## 技术栈

**后端 (backend/):**
- Go 1.21+, Fiber, Uber-FX, Viper, Zerolog
- MacCMS v10 API 客户端
- Swagger 文档

**前端 (frontend/):**
- SvelteKit 2.0 + TypeScript
- Tailwind CSS
- ArtPlayer + HLS.js
- Paraglide (i18n)
- adapter-static (SSG)

## 开发命令

```bash
# 后端
cd backend
go run cmd/server/main.go     # 运行
make swagger                   # 生成 Swagger 文档

# 前端
cd frontend
pnpm dev                       # 开发服务器 (localhost:9595)
pnpm build                     # 构建
```

## 配置

- `backend/configs/config.yaml` - 基础配置
- `backend/configs/config.local.yaml` - 敏感配置（不提交 git）
  - 包含 auth.passwords 和 sources
  - 支持环境变量 CONFIG_LOCAL (base64 编码)

**密码权限控制:**
```yaml
auth:
  enabled: true
  passwords:
    - password: "xxx"
      adult: false  # 不能访问成人源
    - password: "yyy"
      adult: true   # 可以访问成人源
```

## 部署

- **前端:** Cloudflare Pages (路径: frontend, 输出: build)
  - 环境变量: `PUBLIC_API_HOST=https://your-backend.fly.dev`
- **后端:** Fly.io
  - `fly deploy --local-only` (本地构建，包含 config.local.yaml)

## 目录结构

```
backend/
├── cmd/server/main.go      # 入口
├── configs/                # 配置文件
├── internal/
│   ├── config/             # 配置加载
│   ├── handler/            # HTTP 处理器
│   ├── service/            # 业务逻辑
│   ├── source/             # MacCMS API 客户端
│   └── model/              # 数据模型
└── docs/                   # Swagger 文档

frontend/
├── src/
│   ├── routes/             # 页面
│   │   ├── +page.svelte    # 首页
│   │   ├── search/         # 搜索页
│   │   └── player/         # 播放页
│   └── lib/
│       ├── components/     # UI 组件
│       ├── api/client.ts   # API 客户端
│       └── types/          # TypeScript 类型
├── messages/               # i18n (en.json, zh.json)
└── static/                 # 静态资源
```

## 代码规范

- 注释语言：中文或英文，保持一致
- 前端使用 pnpm，不用 npm
- Svelte 5 runes 语法 ($state, $derived, $effect)
- class: 条件样式用模板字符串，不用 `class:xxx={condition}` (有 `/` 字符时会报错)