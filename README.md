# SearchAV

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![SvelteKit](https://img.shields.io/badge/SvelteKit-2.0-FF3E00?style=flat&logo=svelte)](https://kit.svelte.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

**[中文文档](README_CN.md)**

A modern video search aggregation service built with Go backend and SvelteKit frontend.

## Features

- **Multi-source Aggregation** - Search across multiple video sources simultaneously
- **HLS Playback** - Built-in video player with HLS.js support
- **Ad Filtering** - Automatic ad segment removal from video streams
- **Password Protection** - Optional authentication for private deployment
- **Multi-language** - English and Chinese interface support
- **Responsive Design** - Works on desktop and mobile devices
- **Source Switching** - Switch between different video sources for the same content

## Demo

**Live Demo:** https://sav.haose.love/

**Password:** `admin`

## Screenshots

|               Home                |           Search Results            |
|:---------------------------------:|:-----------------------------------:|
| ![Home](images/Snipaste_en_1.png) | ![Search](images/Snipaste_en_2.png) |

|               Player                |               Episodes                |
|:-----------------------------------:|:-------------------------------------:|
| ![Player](images/Snipaste_en_3.png) | ![Episodes](images/Snipaste_en_4.png) |

## Tech Stack

**Backend:**

- Go 1.21+
- Fiber (Web Framework)
- Uber-FX (Dependency Injection)
- Viper (Configuration)
- Zerolog (Logging)
- Swagger (API Documentation)

**Frontend:**

- SvelteKit 2.0
- TypeScript
- Tailwind CSS
- ArtPlayer + HLS.js
- Paraglide (i18n)

## Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- pnpm

### Backend

```bash
cd backend

# Copy and configure
cp configs/config.yaml configs/config.local.yaml
# Edit config.local.yaml to add your sources

# Run
go run cmd/server/main.go
```

### Frontend

```bash
cd frontend

# Install dependencies
pnpm install

# Development
pnpm dev

# Build
pnpm build
```

### Configuration

Create `configs/config.local.yaml` (gitignored) for sensitive configuration:

```yaml
auth:
  enabled: true
  passwords:
    - password: "normal-user-pass"
      adult: false
    - password: "vip-user-pass"
      adult: true  # Can access adult sources

sources:
  - name: "Source Name"
    code: "source_code"
    url: "https://api.example.com/api.php/provide/vod/"
    adult: false
  - name: "Adult Source"
    code: "adult_source"
    url: "https://api.example.com/api.php/provide/vod/"
    adult: true  # Only accessible with adult-enabled password
```

## API Endpoints

| Endpoint      | Method | Description                              |
|---------------|--------|------------------------------------------|
| `/api/search` | GET    | Search videos (`?q=keyword&adult=0\|1`)  |
| `/api/detail` | GET    | Get video details (`?source=xxx&id=xxx`) |
| `/swagger/*`  | GET    | API documentation                        |

## Project Structure

```
SearchAV/
├── backend/
│   ├── cmd/server/         # Entry point
│   ├── configs/            # Configuration files
│   ├── docs/               # Swagger docs
│   └── internal/
│       ├── config/         # Config loading
│       ├── handler/        # HTTP handlers
│       ├── service/        # Business logic
│       ├── source/         # External API client
│       └── model/          # Data models
├── frontend/
│   ├── src/
│   │   ├── routes/         # Pages
│   │   ├── lib/
│   │   │   ├── components/ # UI components
│   │   │   ├── api/        # API client
│   │   │   └── paraglide/  # i18n (generated)
│   │   └── app.css         # Global styles
│   └── messages/           # Translation files
└── README.md
```

---

## Disclaimer

### Service Nature

SearchAV only provides video search services. It does not directly provide, store, or upload any video content. All
search results come from third-party public APIs.

### User Responsibility

When using this service, users must comply with relevant laws and regulations. Users shall not use search results for
infringing activities, such as downloading or distributing unauthorized works.

### Advertisement Warning

All videos come from third-party collection sites. Advertisements appearing in videos are not affiliated with this site.
Do not trust or click any ads in videos to avoid fraud.

### Security Recommendations

It is strongly recommended to set up password protection during deployment to avoid legal risks from public access. Do
not publicly share or distribute instance links.

### Legal Liability

The developers of this project are not responsible for any consequences arising from the use of this project. When using
this project, you must comply with local laws and regulations.

### Usage Scope

This project is for learning and personal use only. Do not use deployed instances for commercial purposes or public
services. Users are solely responsible for any legal issues arising from public sharing.

---

## License

MIT License - see [LICENSE](LICENSE) for details.
