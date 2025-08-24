# TFT Core

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/rsdlab-dk/tft-core)](https://goreportcard.com/report/github.com/rsdlab-dk/tft-core)

Biblioteca Go para desenvolvimento de aplicações Teamfight Tactics (TFT) com integração à API Riot Games.

## Características

- **Riot API Client** - Cliente HTTP otimizado para API Riot Games
- **Rate Limiting** - Sistema inteligente de rate limiting
- **HTTP Handlers** - Handlers prontos para endpoints TFT
- **Logging** - Sistema de logs estruturado com Zap
- **Middleware** - Middlewares para CORS, rate limit, request ID e logging
- **Validação** - Validadores para PUUID, nomes e regiões
- **Error Handling** - Tratamento robusto de erros da API Riot

## Instalação

```bash
go get github.com/rsdlab-dk/tft-core
```

## Uso Rápido

### Cliente Riot API

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/rsdlab-dk/tft-core/riot"
)

func main() {
    client := riot.NewClient("your-riot-api-key")
    
    summoner, err := client.GetSummonerByName(context.Background(), "br1", "PlayerName")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Summoner: %s (Level %d)\n", summoner.Name, summoner.SummonerLevel)
}
```

### HTTP Server com Handlers

```go
package main

import (
    "net/http"
    "os"
    
    tfthttp "github.com/rsdlab-dk/tft-core/http"
    "github.com/rsdlab-dk/tft-core/logger"
    "github.com/rsdlab-dk/tft-core/ratelimit"
    "github.com/rsdlab-dk/tft-core/riot"
)

func main() {
    log, _ := logger.New("development")
    riotClient := riot.NewClient(os.Getenv("RIOT_API_KEY"))
    rateLimiter := ratelimit.NewMemoryLimiter()
    
    mux := http.NewServeMux()
    
    mux.HandleFunc("/summoner/by-name", 
        tfthttp.WithRequestID(log)(
            tfthttp.WithLogging(log)(
                tfthttp.SummonerByNameHandler(riotClient, rateLimiter, log))))
    
    http.ListenAndServe(":8080", mux)
}
```

### Handler Personalizado

```go
func CustomHandler(riotClient *riot.Client, rateLimiter ratelimit.Limiter, log *logger.Logger) http.HandlerFunc {
    return tfthttp.WithCORS(tfthttp.WithRateLimit(rateLimiter, "summoner", log)(func(w http.ResponseWriter, r *http.Request) {
        puuid := r.URL.Query().Get("puuid")
        requestID := logger.GetRequestID(r.Context())
        
        if !tfthttp.ValidatePUUID(puuid, requestID, log, w, r) {
            return
        }
        
        result, err := riotClient.GetSummonerByPUUID(r.Context(), "br1", puuid)
        if err != nil {
            // Error handling automático
            return
        }
        
        tfthttp.WriteJSON(w, result, log, r)
    }))
}
```

## Documentação da API

### Riot Client

```go
client := riot.NewClient("api-key")

// Buscar summoner por nome
summoner, err := client.GetSummonerByName(ctx, "br1", "PlayerName")

// Buscar summoner por PUUID
summoner, err := client.GetSummonerByPUUID(ctx, "br1", "puuid")

// Buscar summoner por ID
summoner, err := client.GetSummonerByID(ctx, "br1", "summonerId")
```

### Rate Limiting

```go
// Memory limiter (desenvolvimento)
limiter := ratelimit.NewMemoryLimiter()

// Verificar se request é permitido
allowed, err := limiter.Allow(ctx, "key", 100, 2*time.Minute)

// Configuração personalizada
config := ratelimit.NewConfig()
config.Rules["custom-endpoint"] = ratelimit.Rule{
    Rate:   500,
    Window: 10 * time.Minute,
}
```

### Logger

```go
log, err := logger.New("development") // ou "production"

log.Info("message", zap.String("key", "value"))
log.Error("error occurred", zap.Error(err))

// Com context (inclui request ID automaticamente)
log.WithContext(ctx).Info("request processed")
```

### HTTP Responses

```go
// Resposta de sucesso
tfthttp.WriteJSON(w, data, log, r)
// Output: {"success": true, "data": {...}}

// Resposta de erro
tfthttp.WriteError(w, "CODE", "message", 400, log, r)
// Output: {"success": false, "error": {"code": "CODE", "message": "message"}}

// Respostas específicas
tfthttp.WriteNotFound(w, "Resource not found", log, r)
tfthttp.WriteBadRequest(w, "Invalid parameter", log, r)
tfthttp.WriteInternalError(w, log, r)
```

### Middlewares

```go
// Request ID automático
handler = tfthttp.WithRequestID(log)(handler)

// Rate limiting
handler = tfthttp.WithRateLimit(rateLimiter, "endpoint", log)(handler)

// CORS
handler = tfthttp.WithCORS(handler)

// Logging de requests
handler = tfthttp.WithLogging(log)(handler)

// Combinar middlewares
handler = tfthttp.WithRequestID(log)(
    tfthttp.WithLogging(log)(
        tfthttp.WithRateLimit(rateLimiter, "endpoint", log)(
            tfthttp.WithCORS(actualHandler))))
```

### Validação

```go
// Validar PUUID
if !tfthttp.ValidatePUUID(puuid, requestID, log, w, r) {
    return // Response automático em caso de erro
}

// Validar nome do summoner
if !tfthttp.ValidateSummonerName(name, requestID, log, w, r) {
    return
}

// Validar região
if !tfthttp.ValidateRegion(region, requestID, log, w, r) {
    return
}
```

## Regiões Suportadas

- `br1` - Brasil
- `eun1` - Europe Nordic & East
- `euw1` - Europe West
- `jp1` - Japan
- `kr` - Korea
- `la1` - Latin America North
- `la2` - Latin America South
- `na1` - North America
- `oc1` - Oceania
- `ph2` - Philippines
- `ru` - Russia
- `sg2` - Singapore
- `th2` - Thailand
- `tr1` - Turkey
- `tw2` - Taiwan
- `vn2` - Vietnam

## Rate Limits Padrão

- **Summoner**: 100 requests / 2 minutos
- **Match**: 100 requests / 2 minutos
- **League**: 100 requests / 2 minutos
- **Match List**: 1000 requests / 10 segundos

## Tratamento de Erros

A biblioteca automaticamente trata erros da API Riot:

- **404** - Recurso não encontrado
- **429** - Rate limit excedido
- **401/403** - Problemas com API key
- **5xx** - Erros do servidor Riot

## Contribuição

1. Fork o projeto
2. Crie uma feature branch (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Add nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

## Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.