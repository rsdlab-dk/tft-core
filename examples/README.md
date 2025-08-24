# 🚀 Instruções para Publicar TFT Core

## 1. Preparação Local

```bash
# Criar diretório e inicializar
mkdir tft-core && cd tft-core
git init

# Copiar todos os arquivos dos artifacts para as respectivas pastas
# Estrutura final:
# tft-core/
# ├── go.mod
# ├── README.md
# ├── LICENSE
# ├── Makefile
# ├── .gitignore
# ├── .github/workflows/ci.yml
# ├── riot/ (todos os arquivos)
# ├── http/ (todos os arquivos) 
# ├── logger/ (todos os arquivos)
# ├── ratelimit/ (todos os arquivos)
# └── examples/
#     ├── basic/main.go
#     └── server/main.go

# Inicializar módulo Go
go mod init github.com/rsdlab-dk/tft-core
go mod tidy
```

## 2. Criar Repositório no GitHub

1. **Acesse:** https://github.com/new
2. **Repository name:** `tft-core`
3. **Description:** `Go library for Teamfight Tactics (TFT) applications with Riot Games API integration`
4. **Public/Private:** Public
5. **Não** marcar "Add README" (já temos)
6. **Criar repositório**

## 3. Primeira Publicação

```bash
# Adicionar remote
git remote add origin https://github.com/rsdlab-dk/tft-core.git

# Commit inicial
git add .
git commit -m "feat: initial release of TFT Core library

- Riot API client with rate limiting
- HTTP handlers and middleware
- Structured logging with Zap  
- Request validation and error handling
- Memory-based rate limiter
- Complete examples and documentation"

# Push inicial
git branch -M main
git push -u origin main
```

## 4. Criar Release/Tag

```bash
# Criar tag v1.0.0
git tag -a v1.0.0 -m "v1.0.0: Initial release

Features:
- Riot Games API client
- HTTP handlers for summoner endpoints
- Rate limiting system
- Structured logging
- Request validation
- CORS and middleware support
- Complete documentation and examples"

# Push da tag
git push origin v1.0.0
```

## 5. Configurar GitHub Release

1. **Acesse:** https://github.com/rsdlab-dk/tft-core/releases/new
2. **Tag version:** `v1.0.0`
3. **Release title:** `v1.0.0 - Initial Release`
4. **Description:**
```markdown
# 🎉 TFT Core v1.0.0

First stable release of TFT Core - a Go library for building Teamfight Tactics applications.

## ✨ Features

- **Riot API Client** - Complete TFT API integration
- **Rate Limiting** - Intelligent rate limiting system  
- **HTTP Handlers** - Ready-to-use endpoints
- **Middleware** - CORS, logging, request ID, rate limiting
- **Validation** - PUUID, summoner name, region validation
- **Error Handling** - Robust error handling for Riot API
- **Examples** - Complete usage examples

## 📦 Installation

```bash
go get github.com/rsdlab-dk/tft-core@v1.0.0
```

## 🚀 Quick Start

```go
import "github.com/rsdlab-dk/tft-core/riot"

client := riot.NewClient("your-api-key")
summoner, err := client.GetSummonerByName(ctx, "br1", "PlayerName")
```

See [README](https://github.com/rsdlab-dk/tft-core#readme) for complete documentation.
```

5. **Marcar:** "Set as the latest release"
6. **Publicar release**

## 6. Verificar Publicação

```bash
# Testar instalação
cd /tmp
mkdir test-tft-core && cd test-tft-core
go mod init test
go get github.com/rsdlab-dk/tft-core@v1.0.0

# Verificar no Go Proxy
# Aguardar ~5 minutos e verificar:
# https://proxy.golang.org/github.com/rsdlab-dk/tft-core/@v/list
# https://pkg.go.dev/github.com/rsdlab-dk/tft-core
```

## 7. Próximos Passos

### Para updates futuros:
```bash
# Fazer mudanças
git add .
git commit -m "feat: add match history support"

# Criar nova tag
git tag -a v1.1.0 -m "v1.1.0: Add match history"
git push origin v1.1.0

# Criar release no GitHub
```

### Configurar badges no README:
- Go Report Card: https://goreportcard.com/report/github.com/rsdlab-dk/tft-core  
- Go Version: Badge automático
- Coverage: Configurar Codecov

### Marketing:
- Post no Reddit r/golang
- Post no Discord/comunidades TFT
- Documentação no pkg.go.dev aparecerá automaticamente

## 🎯 URLs Importantes

- **Repositório:** https://github.com/rsdlab-dk/tft-core
- **Documentação:** https://pkg.go.dev/github.com/rsdlab-dk/tft-core
- **Releases:** https://github.com/rsdlab-dk/tft-core/releases
- **Go Proxy:** https://proxy.golang.org/github.com/rsdlab-dk/tft-core