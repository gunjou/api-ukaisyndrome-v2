package absensi

import (
	"api-ukaisyndrome-v2/pkg/timeutil"
	"fmt"
	"math"
	"time"
)

const (
	maxCheckInDuration = 4 * time.Hour

	pesertaBaseRadiusMeter = 100
	pesertaMaxRadiusMeter  = 300
	defaultGPSAccuracy     = 100

	// Akurasi di dalam ruangan sering terlalu optimis (±7 m padahal meleset puluhan meter).
	minGPSAccuracy = 20
)

func accuracyOrDefault(accuracy *float64) float64 {
	if accuracy == nil {
		return defaultGPSAccuracy
	}
	return *accuracy
}

// Lolos bila lingkaran ketidakpastian kedua titik masih bisa bersinggungan dengan radius dasar.
func pesertaRadius(mentorAccuracy, pesertaAccuracy *float64) float64 {
	return math.Min(
		pesertaMaxRadiusMeter,
		pesertaBaseRadiusMeter+
			math.Max(accuracyOrDefault(mentorAccuracy), minGPSAccuracy)+
			math.Max(accuracyOrDefault(pesertaAccuracy), minGPSAccuracy),
	)
}

func compassDirection(fromLat, fromLon, toLat, toLon float64) string {
	directions := []string{
		"utara", "timur laut", "timur", "tenggara",
		"selatan", "barat daya", "barat", "barat laut",
	}

	fromLatRad := fromLat * math.Pi / 180
	toLatRad := toLat * math.Pi / 180
	dLon := (toLon - fromLon) * math.Pi / 180

	bearing := math.Atan2(
		math.Sin(dLon)*math.Cos(toLatRad),
		math.Cos(fromLatRad)*math.Sin(toLatRad)-
			math.Sin(fromLatRad)*math.Cos(toLatRad)*math.Cos(dLon),
	) * 180 / math.Pi

	return directions[int(math.Round(math.Mod(bearing+360, 360)/45))%8]
}

// Pesan sengaja hanya memuat sisa jarak dan arah agar batas radius tidak terbaca.
func outsideRadiusError(lat, lon, mentorLat, mentorLon, distance, radius float64) error {
	return fmt.Errorf(
		"%w. Bergeser sekitar %.0f meter ke arah %s lalu coba lagi.",
		ErrOutsideRadius,
		math.Ceil(distance-radius),
		compassDirection(lat, lon, mentorLat, mentorLon),
	)
}

func isValidLocation(latitude, longitude float64, accuracy *float64) bool {
	return latitude >= -90 && latitude <= 90 &&
		longitude >= -180 && longitude <= 180 &&
		(accuracy == nil || *accuracy >= 0)
}

func isCheckInExpired(checkInAt time.Time) bool {
	return timeutil.Now().After(
		checkInAt.Add(maxCheckInDuration),
	)
}

func calculateDistance(
	lat1 float64,
	lon1 float64,
	lat2 float64,
	lon2 float64,
) float64 {

	const earthRadius = 6371000

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180

	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*
			math.Cos(lat2Rad)*
			math.Sin(dLon/2)*
			math.Sin(dLon/2)

	c := 2 * math.Atan2(
		math.Sqrt(a),
		math.Sqrt(1-a),
	)

	return earthRadius * c
}