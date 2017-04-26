package affect

// Co-Efficients For Evaluation

// Co-Efficients For Potency

// Co-Efficients For Activity

// AffectiveState - The affective state of an individual throughout an
//                  interaction
type AffectiveState struct {
	Participant          string
	FundamentalSentiment [3]float32
	TransientImpression  [3]float32
	Deflection           [3]float32
}

// CalculateDeflection - Required to Calculate Deflection
func CalculateDeflection(epa1 [3]float32, epa2 [3]float32) float32 {
	return epa1[0] - epa2[0]
}

// CalculateTransient - z
func CalculateTransient(a [3]float32, b [3]float32, o [3]float32) [3]float32 {

	return [3]float32{}
}
