package affect

// Co-Efficients Fpr Evaluation


// Co-Efficients For Potency


// Co-Efficients For Activity




// CalculateDeflection - Required to Calculate Deflection
func CalculateDeflection(epa1 [3]float32, epa2 [3]float32) float32 {
	return epa1[0] - epa2[0]
}

// CalculateTransient - z
func CalculateTransient(a [3]float32, b [3]float32, o [3]float32) [3]float32 {

	return [3]float32{}
}
