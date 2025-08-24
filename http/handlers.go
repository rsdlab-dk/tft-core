package http

import (
	"context"
	"net/http"

	"github.com/rsdlab-dk/tft-core/logger"
	"github.com/rsdlab-dk/tft-core/ratelimit"
	"github.com/rsdlab-dk/tft-core/riot"
	"go.uber.org/zap"
)

func SummonerByNameHandler(riotClient *riot.Client, rateLimiter ratelimit.Limiter, log *logger.Logger) http.HandlerFunc {
	return WithCORS(WithRateLimit(rateLimiter, "summoner", log)(func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		region := r.URL.Query().Get("region")
		if region == "" {
			region = "br1"
		}
		
		requestID := logger.GetRequestID(r.Context())
		
		if !ValidateSummonerName(name, requestID, log, w, r) {
			return
		}
		
		if !ValidateRegion(region, requestID, log, w, r) {
			return
		}

		log.WithContext(r.Context()).Info("summoner request by name",
			zap.String("name", name),
			zap.String("region", region))

		ctx := context.WithValue(r.Context(), "timeout", "10s")
		result, err := riotClient.GetSummonerByName(ctx, region, name)
		if err != nil {
			handleRiotError(err, log, w, r, requestID)
			return
		}

		log.WithContext(r.Context()).Info("summoner request successful",
			zap.String("name", name),
			zap.String("puuid", result.PUUID))

		WriteJSON(w, result, log, r)
	}))
}

func SummonerByPUUIDHandler(riotClient *riot.Client, rateLimiter ratelimit.Limiter, log *logger.Logger) http.HandlerFunc {
	return WithCORS(WithRateLimit(rateLimiter, "summoner", log)(func(w http.ResponseWriter, r *http.Request) {
		puuid := r.URL.Query().Get("puuid")
		region := r.URL.Query().Get("region")
		if region == "" {
			region = "br1"
		}
		
		requestID := logger.GetRequestID(r.Context())
		
		if !ValidatePUUID(puuid, requestID, log, w, r) {
			return
		}
		
		if !ValidateRegion(region, requestID, log, w, r) {
			return
		}

		log.WithContext(r.Context()).Info("summoner request by puuid",
			zap.String("puuid", puuid),
			zap.String("region", region))

		ctx := context.WithValue(r.Context(), "timeout", "10s")
		result, err := riotClient.GetSummonerByPUUID(ctx, region, puuid)
		if err != nil {
			handleRiotError(err, log, w, r, requestID)
			return
		}

		log.WithContext(r.Context()).Info("summoner request successful",
			zap.String("puuid", puuid),
			zap.String("name", result.Name))

		WriteJSON(w, result, log, r)
	}))
}

func handleRiotError(err error, log *logger.Logger, w http.ResponseWriter, r *http.Request, requestID string) {
	if riotErr, ok := err.(*riot.RiotError); ok {
		switch {
		case riotErr.IsNotFound():
			log.WithContext(r.Context()).Warn("summoner not found",
				zap.String("request_id", requestID),
				zap.Error(err))
			WriteNotFound(w, "Summoner not found", log, r)
		case riotErr.IsRateLimited():
			log.WithContext(r.Context()).Warn("rate limited by riot api",
				zap.String("request_id", requestID),
				zap.Error(err))
			WriteError(w, "RATE_LIMITED", "Rate limited by Riot API", http.StatusTooManyRequests, log, r)
		case riotErr.IsUnauthorized() || riotErr.IsForbidden():
			log.WithContext(r.Context()).Error("api key issue",
				zap.String("request_id", requestID),
				zap.Error(err))
			WriteError(w, "API_KEY_ERROR", "API key issue", http.StatusUnauthorized, log, r)
		default:
			log.WithContext(r.Context()).Error("riot api error",
				zap.String("request_id", requestID),
				zap.Error(err))
			WriteInternalError(w, log, r)
		}
		return
	}

	log.WithContext(r.Context()).Error("unexpected error",
		zap.String("request_id", requestID),
		zap.Error(err))
	WriteInternalError(w, log, r)
}