package http

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/rsdlab-dk/tft-core/logger"
	"go.uber.org/zap"
)

var (
	puuidRegex    = regexp.MustCompile(`^[A-Za-z0-9_-]{78}$`)
	gameNameRegex = regexp.MustCompile(`^[\p{L}\p{N}._\- ]{3,16}$`)
	tagLineRegex  = regexp.MustCompile(`^[A-Za-z0-9]{3,5}$`)
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

func ValidateRiotID(gameName, tagLine, requestID string, log *logger.Logger, w http.ResponseWriter, r *http.Request) bool {
	if gameName == "" {
		log.WithContext(r.Context()).Warn("missing gameName parameter",
			zap.String("request_id", requestID))
		WriteBadRequest(w, "GameName parameter is required", log, r)
		return false
	}

	if tagLine == "" {
		log.WithContext(r.Context()).Warn("missing tagLine parameter",
			zap.String("request_id", requestID))
		WriteBadRequest(w, "TagLine parameter is required", log, r)
		return false
	}

	trimmedGameName := strings.TrimSpace(gameName)
	if !gameNameRegex.MatchString(trimmedGameName) {
		log.WithContext(r.Context()).Warn("invalid gameName format",
			zap.String("request_id", requestID),
			zap.String("gameName", gameName))
		WriteBadRequest(w, "Invalid GameName format (3-16 characters)", log, r)
		return false
	}

	trimmedTagLine := strings.TrimSpace(tagLine)
	if !tagLineRegex.MatchString(trimmedTagLine) {
		log.WithContext(r.Context()).Warn("invalid tagLine format",
			zap.String("request_id", requestID),
			zap.String("tagLine", tagLine))
		WriteBadRequest(w, "Invalid TagLine format (3-5 alphanumeric characters)", log, r)
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
		"ru":   true,
		"sg2":  true,
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
