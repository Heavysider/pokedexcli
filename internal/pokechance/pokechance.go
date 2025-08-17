package pokechance

func CalculateCaptureChance(exp int) float64 {
	baseCatchRate := 0.05
	scalingFactor := 100.0
	if exp <= 0 {
		return 0.95
	}

	catchModifier := scalingFactor / (scalingFactor + float64(exp))
	catchChance := baseCatchRate + 0.9*float64(catchModifier)

	return catchChance
}
