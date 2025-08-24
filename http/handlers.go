package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/rsdlab-dk/tft-core/logger"
	"github.com/rsdlab-dk/tft-core/ratelimit"
	"github.com/rsdlab-dk/tft-core/riot"
	"go.uber.org/zap"
)

func SummonerByRiotIDHandler(riotClient *riot.Client, rateLimiter ratelimit.Limiter, log *logger.Logger) http.HandlerFunc {
	return WithCORS(WithRateLimit(rateLimiter, "summoner", log)(func(w http.ResponseWriter, r *http.Request) {
		gameName := r.URL.Query().Get("gameName")
		tagLine := r.URL.Query().Get("tagLine")
		region := r.URL.Query().Get("region")
		if region == "" {
			region = "br1"
		}

		requestID := logger.GetRequestID(r.Context())

		if !ValidateRiotID(gameName, tagLine, requestID, log, w, r) {
			return
		}

		if !ValidateRegion(region, requestID, log, w, r) {
			return
		}

		log.WithContext(r.Context()).Info("summoner request by riot id",
			zap.String("gameName", gameName),
			zap.String("tagLine", tagLine),
			zap.String("region", region))

		ctx := context.WithValue(r.Context(), "timeout", "10s")
		result, err := riotClient.GetSummonerByRiotID(ctx, region, gameName, tagLine)
		if err != nil {
			handleRiotError(err, log, w, r, requestID)
			return
		}

		log.WithContext(r.Context()).Info("summoner request successful",
			zap.String("gameName", gameName),
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

func ChallengerLeagueHandler(riotClient *riot.Client, rateLimiter ratelimit.Limiter, log *logger.Logger) http.HandlerFunc {
	return WithCORS(WithRateLimit(rateLimiter, "league", log)(func(w http.ResponseWriter, r *http.Request) {
		region := r.URL.Query().Get("region")
		if region == "" {
			region = "br1"
		}

		requestID := logger.GetRequestID(r.Context())

		if !ValidateRegion(region, requestID, log, w, r) {
			return
		}

		log.WithContext(r.Context()).Info("challenger league request",
			zap.String("region", region))

		ctx := context.WithValue(r.Context(), "timeout", "10s")
		result, err := riotClient.GetChallengerLeague(ctx, region)
		if err != nil {
			handleRiotError(err, log, w, r, requestID)
			return
		}

		log.WithContext(r.Context()).Info("challenger league request successful",
			zap.String("region", region),
			zap.Int("players", len(result.Entries)))

		WriteJSON(w, result, log, r)
	}))
}

func handleRiotError(err error, log *logger.Logger, w http.ResponseWriter, r *http.Request, requestID string) {
	if riotErr, ok := err.(*riot.RiotError); ok {
		switch {
		case riotErr.IsNotFound():
			log.WithContext(r.Context()).Warn("resource not found",
				zap.String("request_id", requestID),
				zap.Error(err))
			WriteNotFound(w, "Resource not found", log, r)
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
