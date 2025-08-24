package http

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/rsdlab-dk/tft-core/logger"
	"go.uber.org/zap"
)

var (
	puuidRegex = regexp.MustCompile(`^[A-Za-z0-9_-]{78}$`)
	nameRegex  = regexp.MustCompile(`^[0-9A-Za-z ._]{3,16}$`)
)

func ValidatePUUID(puuid, requestID string, log *logger.Logger, w http.ResponseWriter, r *http.Request) bool {
	if puuid == "" {
		log.WithContext(r.Context()).Warn("missing puuid parameter",
			zap.String("request_id", requestID))
		WriteBadRequest(w, "PUUID parameter is required", log, r)
		return false
	}

	if !puuidRegex.MatchString(puuid) {
		log.WithContext(r.Context()).Warn("invalid puuid format",
			zap.String("request_id", requestID),
			zap.String("puuid", puuid))
		WriteBadRequest(w, "Invalid PUUID format", log, r)
		return false
	}

	return true
}

func ValidateSummonerName(name, requestID string, log *logger.Logger, w http.ResponseWriter, r *http.Request) bool {
	if name == "" {
		log.WithContext(r.Context()).Warn("missing summoner name parameter",
			zap.String("request_id", requestID))
		WriteBadRequest(w, "Summoner name parameter is required", log, r)
		return false
	}

	trimmedName := strings.TrimSpace(name)
	if !nameRegex.MatchString(trimmedName) {
		log.WithContext(r.Context()).Warn("invalid summoner name format",
			zap.String("request_id", requestID),
			zap.String("name", name))
		WriteBadRequest(w, "Invalid summoner name format", log, r)
		return false
	}

	return true
}

func ValidateRegion(region, requestID string, log *logger.Logger, w http.ResponseWriter, r *http.Request) bool {
	validRegions := map[string]bool{
		"br1":  true,
		"eun1": true,
		"euw1": true,
		"jp1":  true,
		"kr":   true,
		"la1":  true,
		"la2":  true,
		"na1":  true,
		"oc1":  true,
		"ph2":  true,
		"ru":   true,
		"sg2":  true,
		"th2":  true,
		"tr1":  true,
		"tw2":  true,
		"vn2":  true,
	}

	if region == "" {
		region = "br1"
		return true
	}

	if !validRegions[strings.ToLower(region)] {
		log.WithContext(r.Context()).Warn("invalid region",
			zap.String("request_id", requestID),
			zap.String("region", region))
		WriteBadRequest(w, "Invalid region", log, r)
		return false
	}

	return true
}