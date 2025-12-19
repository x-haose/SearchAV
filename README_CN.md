# SearchAV

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![SvelteKit](https://img.shields.io/badge/SvelteKit-2.0-FF3E00?style=flat&logo=svelte)](https://kit.svelte.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

**[English](README.md)**

一个现代化的视频搜索聚合服务，采用 Go 后端 + SvelteKit 前端架构。

## 功能特性

- **多源聚合** - 同时搜索多个视频源
- **HLS 播放** - 内置视频播放器，支持 HLS.js
- **广告过滤** - 自动过滤视频流中的广告片段
- **密码保护** - 支持私有部署的身份验证
- **多语言** - 支持中英文界面
- **响应式设计** - 适配桌面和移动设备
- **源切换** - 同一内容支持多源切换

## 在线演示

**演示地址：** https://sav.haose.love/

**密码：** `admin`

## 截图预览

|               首页                |              搜索结果               |
|:-------------------------------:|:-------------------------------:|
| ![首页](images/Snipaste_zh_1.png) | ![搜索](images/Snipaste_zh_2.png) |

|               播放器                |               选集                |
|:--------------------------------:|:-------------------------------:|
| ![播放器](images/Snipaste_zh_3.png) | ![选集](images/Snipaste_zh_4.png) |

## 技术栈

**后端：**

- Go 1.21+
- Fiber (Web 框架)
- Uber-FX (依赖注入)
- Viper (配置管理)
- Zerolog (日志)
- Swagger (API 文档)

**前端：**

- SvelteKit 2.0
- TypeScript
- Tailwind CSS
- ArtPlayer + HLS.js
- Paraglide (国际化)

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- pnpm

### 后端

```bash
cd backend

# 复制并配置
cp configs/config.yaml configs/config.local.yaml
# 编辑 config.local.yaml 添加你的源

# 运行
go run cmd/server/main.go
```

### 前端

```bash
cd frontend

# 安装依赖
pnpm install

# 开发模式
pnpm dev

# 构建
pnpm build
```

### 配置说明

创建 `configs/config.local.yaml`（已被 gitignore）用于敏感配置：

```yaml
auth:
  enabled: true
  passwords:
    - password: "普通用户密码"
      adult: false
    - password: "VIP用户密码"
      adult: true  # 可访问成人源

sources:
  - name: "源名称"
    code: "source_code"
    url: "https://api.example.com/api.php/provide/vod/"
    adult: false
  - name: "成人源"
    code: "adult_source"
    url: "https://api.example.com/api.php/provide/vod/"
    adult: true  # 仅限 adult: true 的密码访问
```

## API 接口

| 接口            | 方法  | 说明                            |
|---------------|-----|-------------------------------|
| `/api/search` | GET | 搜索视频 (`?q=关键词&adult=0\|1`)    |
| `/api/detail` | GET | 获取视频详情 (`?source=xxx&id=xxx`) |
| `/swagger/*`  | GET | API 文档                        |

## 项目结构

```
SearchAV/
├── backend/
│   ├── cmd/server/         # 入口
│   ├── configs/            # 配置文件
│   ├── docs/               # Swagger 文档
│   └── internal/
│       ├── config/         # 配置加载
│       ├── handler/        # HTTP 处理器
│       ├── service/        # 业务逻辑
│       ├── source/         # 外部 API 客户端
│       └── model/          # 数据模型
├── frontend/
│   ├── src/
│   │   ├── routes/         # 页面
│   │   ├── lib/
│   │   │   ├── components/ # UI 组件
│   │   │   ├── api/        # API 客户端
│   │   │   └── paraglide/  # 国际化（生成）
│   │   └── app.css         # 全局样式
│   └── messages/           # 翻译文件
└── README.md
```

---

## 免责声明

### 服务性质

SearchAV 仅提供视频搜索服务，不直接提供、存储或上传任何视频内容。所有搜索结果均来自第三方公开接口。

### 用户责任

用户在使用本站服务时，须遵守相关法律法规，不得利用搜索结果从事侵权行为，如下载、传播未经授权的作品等。

### 广告风险提示

本站所有视频均来自第三方采集站，视频中出现的广告与本站无关，请勿相信或点击视频中的任何广告内容，谨防上当受骗。

### 安全建议

强烈建议在部署时设置密码保护，避免公开访问可能带来的法律风险。请勿将实例链接公开分享或传播。

### 法律责任

本项目开发者不对使用本项目产生的任何后果负责。使用本项目时，您必须遵守当地的法律法规。

### 使用范围

本项目仅供学习和个人使用，请勿将部署的实例用于商业用途或公开服务。如因公开分享导致的任何法律问题，用户需自行承担责任。

---

## 许可证

MIT License - 详见 [LICENSE](LICENSE)
