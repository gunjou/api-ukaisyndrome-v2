package timeutil

import "time"

var jakartaLocation *time.Location

func init() {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic("failed to load Asia/Jakarta timezone")
	}
	jakartaLocation = loc
}

// ==========================
// GET CURRENT TIME (WIB)
// ==========================
func Now() time.Time {
	return time.Now().In(jakartaLocation)
}

// ==========================
// FORMAT TIME (OPTIONAL)
// ==========================
func NowString() string {
	return Now().Format(time.RFC3339)
}

// ==========================
// NORMALIZE DB TIME (WIB)
// ==========================
func NormalizeDBTime(t time.Time) time.Time {
	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Nanosecond(),
		jakartaLocation,
	)
}