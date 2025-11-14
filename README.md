# go-shortly 🌐

A production-grade URL shortener built in Go.

[![Deploy](https://img.shields.io/badge/Deploy-Render-blue?logo=render)](https://render.com)
[![Go](https://img.shields.io/badge/Go-1.22-00ADD8?logo=go)](https://golang.org)

## ✨ Features

- 🔗 Create short URLs via API or CLI
- ⚡ Redis caching for hot URLs
- 🔒 Rate limiting (10 reqs/sec per IP)
- 📦 CLI tool: `shortly create https://...`
- 📚 OpenAPI docs at `/swagger`

## 🚀 Try It

```bash
# CLI
go install github.com/yourname/go-shortly/cmd/shortly@latest
shortly create https://github.com

# API
curl -X POST https://go-shortly.onrender.com/shorten -d '{"url":"https://render.com"}'
```
