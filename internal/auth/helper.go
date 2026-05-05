package auth

import (
	"strconv"
	"time"
)

// ==========================
// TTL
// ==========================
func getAccessTTL(platform string) time.Duration {
	if platform == "dev" {
        return 3 * time.Hour
    }
    if platform == "mobile" {
        return 8760 * time.Hour // 1 tahun
    }
    return 15 * time.Minute
}

// ==========================
// SESSION KEY
// ==========================
func (s *Service) buildSessionKey(userID int, role string, platform string, jti string) string {

	if platform == "web" || platform == "mobile" {
		return "session:" + role + ":" + platform + ":" + strconv.Itoa(userID)
	}

	return "session:" + role + ":" + platform + ":" + strconv.Itoa(userID) + ":" + jti
}

// ==========================
// REFRESH KEY
// ==========================
func buildRefreshKey(userID int, platform string) string {
	return "session:refresh:" + platform + ":" + strconv.Itoa(userID)
}