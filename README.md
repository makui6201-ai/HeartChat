# HeartChat 💬❤️

A WeChat mini-program **virtual companion** powered by the DeepSeek API, built with a Go backend.

---

## ✨ Features

| Feature | Detail |
|---|---|
| 🤖 Two characters | 浩然 (温柔男友) · 米娜 (可爱女友) |
| ❤️ Love score system | 4 levels: 普通朋友 → 熟悉 → 亲密 → 暧昧 |
| 💬 Typewriter effect | AI replies appear character by character |
| ⏳ Typing indicator | Animated dots while waiting for AI |
| 💾 Persistent state | Love score saved to device local storage |
| 🚀 Go backend | Stateless REST API, DeepSeek integration |

---

## 📁 Project Structure

```
├── server/                 # Go backend
│   ├── main.go
│   ├── handlers/chat.go    # POST /api/chat, POST /api/greeting
│   ├── services/deepseek.go
│   ├── Dockerfile
│   └── go.mod
│
└── miniprogram/            # WeChat Mini-Program
    ├── app.js / app.json / app.wxss
    ├── pages/
    │   ├── index/          # Character selection page
    │   └── chat/           # Chat page
    └── utils/
        ├── api.js          # API client helpers
        └── config.js       # Character data & server URL
```

---

## 🚀 Quick Start

### 1 · Backend server

```bash
cd server

# Set your DeepSeek API key (and optionally restrict CORS origins)
cp ../.env.example .env
# Edit .env:
#   DEEPSEEK_API_KEY=sk-...
#   CORS_ORIGIN=https://servicewechat.com   ← production only; omit for dev (allows *)

# Run locally
export $(cat ../.env | xargs)
go mod tidy
go run main.go
```

Or with Docker:

```bash
cd server
docker build -t heartchat-server .
docker run -p 8080:8080 -e DEEPSEEK_API_KEY=sk-... heartchat-server
```

### 2 · WeChat Mini-Program

1. Download [WeChat Developer Tools](https://developers.weixin.qq.com/miniprogram/dev/devtools/download.html)
2. **Import** the `miniprogram/` directory as a new project
3. In `miniprogram/utils/config.js` set `BASE_URL` to your server address
4. In Developer Tools → **Details → Local settings** tick **不校验合法域名** for local testing
5. Click **Preview** or **Upload**

> **Production note:** The character portrait images are served from
> `lh3.googleusercontent.com`. Whitelist this domain in your WeChat
> mini-program server settings, or replace the URLs in `config.js` with
> images hosted on your own CDN.

---

## 🔌 API Reference

### `POST /api/chat`

| Field | Type | Description |
|---|---|---|
| `role` | string | `boyfriend` \| `girlfriend` |
| `love` | int | current love score (≥ 0) |
| `message` | string | user's message (max 500 chars) |
| `history` | array | last N `{role, content}` turns |

**Response:**
```json
{ "reply": "辛苦啦，要不要我陪你聊会儿？" }
```

### `POST /api/greeting`

| Field | Type | Description |
|---|---|---|
| `role` | string | `boyfriend` \| `girlfriend` |
| `love` | int | current love score |

**Response:**
```json
{ "greeting": "你来啦～今天在干嘛呀？" }
```

### `GET /health`

Returns `{"status":"ok"}`.

---

## ❤️ Love Score System

| Score | Level |
|---|---|
| 0 – 9 | 普通朋友 |
| 10 – 29 | 熟悉 |
| 30 – 59 | 亲密 |
| 60 + | 暧昧 |

Each sent message awards **+1** love. The AI's tone and greeting
become progressively warmer as the score rises.