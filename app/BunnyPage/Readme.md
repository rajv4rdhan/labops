# BunnyPage - A Quiet Meadow (static Go server)

A tiny static site served by a minimal Go web server. A cute, animated
meadow scene rendered entirely with HTML/CSS — no JavaScript, no backend
calls. It exists to demonstrate a minimal `scratch`-based Docker image.

## 🐰 What is BunnyPage?

An HTTP file server written in Go that serves a single charming `index.html`
featuring a peaceful meadow, drifting clouds and a cute bunny. There is no
API and no database: the page is fully self-contained.

## 📁 Project structure

```text
BunnyPage/
├── main.go       # Go static file server
├── index.html    # Self-contained animated meadow page
├── Dockerfile    # Multi-stage build -> scratch runtime
└── Readme.md     # This file
```

## 🚀 Run locally

```bash
go run .
# open http://localhost:8080
```

## 🐳 Docker

```bash
docker build -t bunnypage .
docker run -p 8080:8080 bunnypage
```

The runtime image is built `FROM scratch` and ships only the statically
compiled server plus `index.html`.

## 🔧 Environment

| Env    | Default | Description      |
| ------ | ------- | ---------------- |
| `PORT` | `8080`  | HTTP listen port |
