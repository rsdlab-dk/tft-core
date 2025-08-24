package main

import (
	"context"
	"net/http"
	"os"
	"time"

	tfthttp "github.com/rsdlab-dk/tft-core/http"
	"github.com/rsdlab-dk/tft-core/logger"
	"github.com/rsdlab-dk/tft-core/ratelimit"
	"github.com/rsdlab-dk/tft-core/riot"
	"go.uber.org/zap"
)

func main() {
	log, err := logger.New("development")
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	apiKey := os.Getenv("RIOT_API_KEY")
	if apiKey == "" {
		log.Fatal("RIOT_API_KEY environment variable is required")
	}

	riotClient := riot.NewClient(apiKey)
	rateLimiter := ratelimit.NewMemoryLimiter()

	mux := http.NewServeMux()

	mux.HandleFunc("/summoner/by-name", 
		tfthttp.WithRequestID(log)(
			tfthttp.WithLogging(log)(
				tfthttp.SummonerByNameHandler(riotClient, rateLimiter, log))))

	mux.HandleFunc("/summoner/by-puuid", 
		tfthttp.WithRequestID(log)(
			tfthttp.WithLogging(log)(
				tfthttp.SummonerByPUUIDHandler(riotClient, rateLimiter, log))))

	mux.HandleFunc("/player/stats", CustomPlayerHandler(riotClient, rateLimiter, log))

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		tfthttp.WriteJSON(w, map[string]interface{}{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"version":   "1.0.0",
		}, log, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Info("🚀 TFT Core Server iniciado", 
		zap.String("port", port),
		zap.String("endpoints", "/summoner/by-name, /summoner/by-puuid, /player/stats, /health"))
	
	log.Info("📖 Exemplos de uso:",
		zap.String("summoner-name", "http://localhost:"+port+"/summoner/by-name?name=Hide%20on%20bush&region=kr"),
		zap.String("summoner-puuid", "http://localhost:"+port+"/summoner/by-puuid?puuid=PUUID_HERE&region=kr"),
		zap.String("health", "http://localhost:"+port+"/health"))

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal("Falha ao iniciar servidor", zap.Error(err))
	}
}

func CustomPlayerHandler(riotClient *riot.Client, rateLimiter ratelimit.Limiter, log *logger.Logger) http.HandlerFunc {
	return tfthttp.WithRequestID(log)(
		tfthttp.WithLogging(log)(
			tfthttp.WithCORS(tfthttp.WithRateLimit(rateLimiter, "summoner", log)(func(w http.ResponseWriter, r *http.Request) {
				name := r.URL.Query().Get("name")
				region := r.URL.Query().Get("region")
				if region == "" {
					region = "br1"
				}
				
				requestID := logger.GetRequestID(r.Context())
				
				if !tfthttp.ValidateSummonerName(name, requestID, log, w, r) {
					return
				}
				
				if !tfthttp.ValidateRegion(region, requestID, log, w, r) {
					return
				}

				log.WithContext(r.Context()).Info("player stats request",
					zap.String("name", name),
					zap.String("region", region))

				ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
				defer cancel()

				summoner, err := riotClient.GetSummonerByName(ctx, region, name)
				if err != nil {
					if riotErr, ok := err.(*riot.RiotError); ok {
						switch {
						case riotErr.IsNotFound():
							tfthttp.WriteNotFound(w, "Player not found", log, r)
						case riotErr.IsRateLimited():
							tfthttp.WriteError(w, "RATE_LIMITED", "Rate limited by Riot API", http.StatusTooManyRequests, log, r)
						default:
							log.WithContext(r.Context()).Error("riot api error", zap.Error(err))
							tfthttp.WriteInternalError(w, log, r)
						}
						return
					}
					log.WithContext(r.Context()).Error("unexpected error", zap.Error(err))
					tfthttp.WriteInternalError(w, log, r)
					return
				}

				playerStats := map[string]interface{}{
					"summoner":     summoner,
					"region":       region,
					"last_updated": time.Now().Unix(),
					"stats": map[string]interface{}{
						"rank":        "Master",
						"lp":          1337,
						"games":       42,
						"avg_placement": 3.2,
						"top4_rate":   0.78,
					},
				}

				log.WithContext(r.Context()).Info("player stats successful",
					zap.String("name", summoner.Name),
					zap.String("puuid", summoner.PUUID))

				tfthttp.WriteJSON(w, playerStats, log, r)
			}))))
}