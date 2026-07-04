package score

type ScoreRange struct {
	Label string
	Min   float64
	Max   float64
}

var Distribution = []ScoreRange{
	{
		Label: "0-20",
		Min:   0,
		Max:   20,
	},
	{
		Label: "21-40",
		Min:   21,
		Max:   40,
	},
	{
		Label: "41-60",
		Min:   41,
		Max:   60,
	},
	{
		Label: "61-80",
		Min:   61,
		Max:   80,
	},
	{
		Label: "81-100",
		Min:   81,
		Max:   100,
	},
}

// =================================================
// GET LABELS
// =================================================
func Labels() []string {

	labels := make([]string, 0, len(Distribution))

	for _, r := range Distribution {
		labels = append(labels, r.Label)
	}

	return labels
}

// =================================================
// EMPTY DISTRIBUTION
// =================================================
func EmptyDistribution() map[string]int {

	result := make(map[string]int)

	for _, r := range Distribution {
		result[r.Label] = 0
	}

	return result
}

// =================================================
// GET LABEL FROM SCORE
// =================================================
func GetLabel(value float64) string {

	for _, r := range Distribution {

		if value >= r.Min && value <= r.Max {
			return r.Label
		}

	}

	return "unknown"
}