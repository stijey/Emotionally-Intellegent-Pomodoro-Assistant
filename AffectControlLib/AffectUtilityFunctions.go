package affect

// Co-Efficients For Evaluation

// Co-Efficients For Potency

// Co-Efficients For Activity

// AffectiveState - The affective state of an individual throughout an
//                  interaction

func Deflection() int {
	return 1
}

func PomodoroTime() int {
	return 4
}

func BreakTime() int {
	return 8
}

type AffectiveState struct {
	Participant          string
	FundamentalSentiment [3]float32
	TransientImpression  [3]float32
	Deflection           [3]float32
}


func CalculateTransient(behaviour [3]float32) [3]float32 {
	oldE := behaviour[0]
	oldP := behaviour[1]
	oldA := behaviour[2]

	newE := (oldE * EV_A_A_e) + (oldE * EV_A_A_p) + (oldE * EV_A_A_a) +
		(oldE * EV_A_B_e) + (oldE * EV_A_B_p) + (oldE * EV_A_B_a) +
		(oldE * EV_A_O_e) + (oldE * EV_A_O_p) + (oldE * EV_A_O_a) +
		(oldE * EV_AB_A_ea) + (oldE * EV_BO_A_ee) + (oldE * EV_AB_A_pp) +
		(oldE * EV_BO_A_pp) + (oldE * EV_AB_A_aa) + (oldE * EV_AB_A_pe) +
		(oldE * EV_AO_A_pa) + (oldE * EV_BO_A_ep) + (oldE * EV_BO_A_pe) +
		(oldE * EV_BO_A_pa) + (oldE * EV_BO_A_pe) + (oldE * EV_BO_A_pa) +
		(oldE * EV_BO_A_ae) + (oldE * EV_BO_A_ap) + (oldE * EV_ABO_A_eee) +
		(oldE * EV_ABO_A_ppp) + (oldE * EV_ABO_A_aaa) + (oldE * EV_ABO_A_epp) +
		(oldE * EV_ABO_A_ppa) + EV_A_A_constant

	newP := (oldP * EV_A_A_e) + (oldP * EV_A_A_p) + (oldP * EV_A_A_a) +
		(oldP * EV_A_B_e) + (oldP * EV_A_B_p) + (oldP * EV_A_B_a) +
		(oldP * EV_A_O_e) + (oldP * EV_A_O_p) + (oldP * EV_A_O_a) +
		(oldP * EV_AB_A_ea) + (oldP * EV_BO_A_ee) + (oldP * EV_AB_A_pp) +
		(oldP * EV_BO_A_pp) + (oldP * EV_AB_A_aa) + (oldP * EV_AB_A_pe) +
		(oldP * EV_AO_A_pa) + (oldP * EV_BO_A_ep) + (oldP * EV_BO_A_pe) +
		(oldP * EV_BO_A_pa) + (oldP * EV_BO_A_pe) + (oldP * EV_BO_A_pa) +
		(oldP * EV_BO_A_ae) + (oldP * EV_BO_A_ap) + (oldP * EV_ABO_A_eee) +
		(oldP * EV_ABO_A_ppp) + (oldP * EV_ABO_A_aaa) + (oldP * EV_ABO_A_epp) +
		(oldP * EV_ABO_A_ppa) + EV_A_A_constant

	newA := (oldA * EV_A_A_e) + (oldA * EV_A_A_p) + (oldA * EV_A_A_a) +
		(oldA * EV_A_B_e) + (oldA * EV_A_B_p) + (oldA * EV_A_B_a) +
		(oldA * EV_A_O_e) + (oldA * EV_A_O_p) + (oldA * EV_A_O_a) +
		(oldA * EV_AB_A_ea) + (oldA * EV_BO_A_ee) + (oldA * EV_AB_A_pp) +
		(oldA * EV_BO_A_pp) + (oldA * EV_AB_A_aa) + (oldA * EV_AB_A_pe) +
		(oldA * EV_AO_A_pa) + (oldA * EV_BO_A_ep) + (oldA * EV_BO_A_pe) +
		(oldA * EV_BO_A_pa) + (oldA * EV_BO_A_pe) + (oldA * EV_BO_A_pa) +
		(oldA * EV_BO_A_ae) + (oldA * EV_BO_A_ap) + (oldA * EV_ABO_A_eee) +
		(oldA * EV_ABO_A_ppp) + (oldA * EV_ABO_A_aaa) + (oldA * EV_ABO_A_epp) +
		(oldA * EV_ABO_B_ppa) + EV_A_A_constant

	return [3]float32{newE, newP, newA}
}



// // CalculateDeflection - Required to Calculate Deflection
// func CalculateDeflection(actor [3]float32, behaviour [3]float32, object [3]float32) {
// 	newActorEPA := calculateActorTransient(actor)
// 	//calculateBehaviourTransient()
// 	//calculateObjectTransient()
// }
//
// // CalculateTransient - z
// func CalculateTransient(a [3]float32, b [3]float32, o [3]float32) [3]float32 {
//
// 	return [3]float32{}
// }
