package criteria

func almostEqual(a, b, tolerance float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= tolerance
}
