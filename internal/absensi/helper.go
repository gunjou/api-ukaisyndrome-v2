package absensi

import (
	"api-ukaisyndrome-v2/pkg/timeutil"
	"math"
	"time"
)

const maxCheckInDuration = 4 * time.Hour

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