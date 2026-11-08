package absensi

import "testing"

func ptr(v float64) *float64 { return &v }

func TestPesertaRadius(t *testing.T) {
	cases := []struct {
		mentor, peserta *float64
		want            float64
	}{
		{ptr(10), ptr(10), 120},
		{ptr(50), ptr(50), 200},
		{ptr(100), ptr(100), 300},
		{ptr(300), ptr(300), 300},
		{nil, ptr(10), 210},
		{nil, nil, 300},
	}

	for _, c := range cases {
		if got := pesertaRadius(c.mentor, c.peserta); got != c.want {
			t.Errorf("pesertaRadius = %v, want %v", got, c.want)
		}
	}
}

func TestCompassDirection(t *testing.T) {
	const lat, lon, d = -6.2, 106.8, 0.001

	cases := []struct {
		toLat, toLon float64
		want         string
	}{
		{lat + d, lon, "utara"},
		{lat + d, lon + d, "timur laut"},
		{lat, lon + d, "timur"},
		{lat - d, lon + d, "tenggara"},
		{lat - d, lon, "selatan"},
		{lat - d, lon - d, "barat daya"},
		{lat, lon - d, "barat"},
		{lat + d, lon - d, "barat laut"},
	}

	for _, c := range cases {
		if got := compassDirection(lat, lon, c.toLat, c.toLon); got != c.want {
			t.Errorf("compassDirection = %v, want %v", got, c.want)
		}
	}
}

func TestOutsideRadiusMessageHidesLimit(t *testing.T) {
	err := outsideRadiusError(-6.2, 106.8, -6.2, 106.801, 122.4, 120)

	want := "Anda berada di luar area kelas. Bergeser sekitar 3 meter ke arah timur lalu coba lagi."
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestIsValidLocation(t *testing.T) {
	cases := []struct {
		lat, lon float64
		accuracy *float64
		want     bool
	}{
		{-6.2, 106.8, ptr(10), true},
		{-6.2, 106.8, nil, true},
		{91, 106.8, ptr(10), false},
		{-6.2, 181, ptr(10), false},
		{-6.2, 106.8, ptr(-1), false},
	}

	for _, c := range cases {
		if got := isValidLocation(c.lat, c.lon, c.accuracy); got != c.want {
			t.Errorf("isValidLocation(%v, %v) = %v, want %v", c.lat, c.lon, got, c.want)
		}
	}
}
