package affect

import "math"
import "strings"
import "io/ioutil"
import "strconv"

type AffectiveState struct {
	Actor                [3]float64
	Behaviour            [3]float64
	Object               [3]float64
	FundamentalSentiment [3]float64
	TransientImpression  [3]float64
	Deflection           float64
	Behaviours           map[string][3]float64
}

func MakeAffectiveState() AffectiveState {
	actor := LoadIdentities("student")
	object := LoadIdentities("secretary")
	behaviours := LoadBehaviours()

	a := AffectiveState{Actor: actor,
		Behaviour:            [3]float64{1, 1, 1},
		Object:               object,
		FundamentalSentiment: [3]float64{1, 1, 1},
		TransientImpression:  [3]float64{1, 1, 1},
		Deflection:           0,
		Behaviours:           behaviours}

	return a
}

func (a *AffectiveState) PropegateForward(behaviour [3]float64) {
	// Re-calculate the values
	a.TransientImpression = CalculateTransient(a.Actor, behaviour, a.Object)
	a.Deflection = CalculateDeflection(a.TransientImpression, a.Actor,
		behaviour, a.Behaviour, a.Object, a.Object)
	a.FundamentalSentiment = CalculateTransient(a.Actor, behaviour, a.Object)

	//Switch the users for the next turn
	temp := a.Actor
	a.Actor = a.Object
	a.Object = temp

}

func (a *AffectiveState) Respond() {
	if a.Deflection > 21 {
		a.PropegateForward(a.Behaviours["compliment"])
	} else if a.Deflection > 10 {
		a.PropegateForward(a.Behaviours["counsel"])
	} else {
		a.PropegateForward(a.Behaviours["pursue"])
	}
}

func CalculateTransient(actor [3]float64, behaviour [3]float64,
	object [3]float64) [3]float64 {

	oldActorE := actor[0]
	oldActorP := actor[1]
	oldActorA := actor[2]

	behaviourE := behaviour[0]
	behaviourP := behaviour[1]
	behaviourA := behaviour[2]

	objectE := object[0]
	objectP := object[1]
	objectA := object[2]

	newE := ((oldActorE * EV_A_A_e) - (oldActorP * EV_A_A_p) - (oldActorA * EV_A_A_a)) +
		((behaviourE * EV_A_B_e) - (behaviourP * EV_A_B_p) - (behaviourA * EV_A_B_a)) +
		((objectE * EV_A_O_e) - (objectP * EV_A_O_p) - (objectA * EV_A_O_a)) +

		((oldActorE * behaviourE) * EV_AB_A_ea) + ((behaviourE * objectE) * EV_BO_A_ee) + ((oldActorP * behaviourP) * EV_AB_A_pp) + ((behaviourP * objectP) * EV_BO_A_pp) +

		(((oldActorA * behaviourA) * EV_AB_A_aa) - ((oldActorE * behaviourP) * EV_AB_A_ep) - ((oldActorE * behaviourA) * EV_AB_A_ea)) +
		(((oldActorP * behaviourE) * EV_AB_A_pe) - ((oldActorP * objectA) * EV_BO_A_pa) - ((behaviourE * objectP) * EV_BO_A_ep) - ((behaviourP * objectA) * EV_BO_A_pa)) +

		((behaviourA * objectE) * EV_BO_A_ae) + ((behaviourA * objectP) * EV_BO_A_ap) +
		(((oldActorE * behaviourE * objectE) * EV_ABO_A_eee) - ((oldActorP * behaviourP * objectP) * EV_ABO_A_ppp)) +
		((oldActorA * behaviourA * objectA) * EV_ABO_A_aaa) + ((oldActorE * behaviourP * objectP) * EV_ABO_A_epp) + ((oldActorP * behaviourP * objectA) * EV_ABO_A_ppa) +
		EV_A_A_constant

	newP := ((oldActorE * PO_A_A_e) - (oldActorP * PO_A_A_p) - (oldActorA * PO_A_A_a)) +
		((behaviourE * PO_A_B_e) - (behaviourP * PO_A_B_p) - (behaviourA * PO_A_B_a)) +
		((objectE * PO_A_O_e) - (objectP * PO_A_O_p) - (objectA * PO_A_O_a)) +

		((oldActorE * behaviourE) * PO_AB_A_ea) + ((behaviourE * objectE) * PO_BO_A_ee) + ((oldActorP * behaviourP) * PO_AB_A_pp) + ((behaviourP * objectP) * PO_BO_A_pp) +

		(((oldActorA * behaviourA) * PO_AB_A_aa) - ((oldActorE * behaviourP) * PO_AB_A_ep) - ((oldActorE * behaviourA) * PO_AB_A_ea)) +
		(((oldActorP * behaviourE) * PO_AB_A_pe) - ((oldActorP * objectA) * PO_BO_A_pa) - ((behaviourE * objectP) * PO_BO_A_ep) - ((behaviourP * objectA) * PO_BO_A_pa)) +

		((behaviourA * objectE) * PO_BO_A_ae) + ((behaviourA * objectP) * PO_BO_A_ap) +
		(((oldActorE * behaviourE * objectE) * PO_ABO_A_eee) - ((oldActorP * behaviourP * objectP) * PO_ABO_A_ppp)) +
		((oldActorA * behaviourA * objectA) * PO_ABO_A_aaa) + ((oldActorE * behaviourP * objectP) * PO_ABO_A_epp) + ((oldActorP * behaviourP * objectA) * PO_ABO_A_ppa) +
		PO_A_B_constant

	newA := ((oldActorE * AC_A_A_e) - (oldActorP * AC_A_A_p) - (oldActorA * AC_A_A_a)) +
		((behaviourE * AC_A_B_e) - (behaviourP * AC_A_B_p) - (behaviourA * AC_A_B_a)) +
		((objectE * AC_A_O_e) - (objectP * AC_A_O_p) - (objectA * AC_A_O_a)) +

		((oldActorE * behaviourE) * AC_AB_A_ea) + ((behaviourE * objectE) * AC_BO_A_ee) + ((oldActorP * behaviourP) * AC_AB_A_pp) + ((behaviourP * objectP) * AC_BO_A_pp) +

		(((oldActorA * behaviourA) * AC_AB_A_aa) - ((oldActorE * behaviourP) * AC_AB_A_ep) - ((oldActorE * behaviourA) * AC_AB_A_ea)) +
		(((oldActorP * behaviourE) * AC_AB_A_pe) - ((oldActorP * objectA) * AC_BO_A_pa) - ((behaviourE * objectP) * AC_BO_A_ep) - ((behaviourP * objectA) * AC_BO_A_pa)) +

		((behaviourA * objectE) * AC_BO_A_ae) + ((behaviourA * objectP) * AC_BO_A_ap) +
		(((oldActorE * behaviourE * objectE) * AC_ABO_A_eee) - ((oldActorP * behaviourP * objectP) * AC_ABO_A_ppp)) +
		((oldActorA * behaviourA * objectA) * AC_ABO_A_aaa) + ((oldActorE * behaviourP * objectP) * AC_ABO_A_epp) + ((oldActorP * behaviourP * objectA) * AC_ABO_A_ppa) +
		AC_A_O_constant

	return [3]float64{newE, newP, newA}
}

func CalculateDeflection(actor [3]float64, oldActor [3]float64,
	oldBehaviour [3]float64, behaviour [3]float64,
	oldObject [3]float64, object [3]float64) float64 {

	deflection := math.Pow(actor[0]-oldActor[0], 2) +
		math.Pow(actor[1]-oldActor[1], 2) +
		math.Pow(actor[2]-oldActor[2], 2) +
		math.Pow(behaviour[0]-oldBehaviour[0], 2) +
		math.Pow(behaviour[1]-oldBehaviour[1], 2) +
		math.Pow(behaviour[2]-oldBehaviour[2], 2) +
		math.Pow(object[0]-oldObject[0], 2) +
		math.Pow(object[1]-oldObject[1], 2) +
		math.Pow(object[2]-oldObject[2], 2)

	return deflection
}

func LoadIdentities(ident string) [3]float64 {
	beh := make(map[string][3]float64)
	b, _ := ioutil.ReadFile("AffectControlLib/affect-data/identities.dat")
	str := string(b)
	buffer := []string(strings.Fields(str))

	for i := 0; i < len(buffer); i++ {
		data1, _ := strconv.ParseFloat(buffer[i+1], 64)
		data2, _ := strconv.ParseFloat(buffer[i+2], 64)
		data3, _ := strconv.ParseFloat(buffer[i+3], 64)
		beh[buffer[i]] = [3]float64{data1, data2, data3}
		i += 3
	}
	return beh[ident]
}

func LoadBehaviours() map[string][3]float64 {
	beh := make(map[string][3]float64)
	b, _ := ioutil.ReadFile("AffectControlLib/affect-data/behaviours.dat")
	str := string(b)
	buffer := []string(strings.Fields(str))

	for i := 0; i < len(buffer); i++ {
		data1, _ := strconv.ParseFloat(buffer[i+1], 64)
		data2, _ := strconv.ParseFloat(buffer[i+2], 64)
		data3, _ := strconv.ParseFloat(buffer[i+3], 64)
		beh[buffer[i]] = [3]float64{data1, data2, data3}
		i += 6
	}
	return beh
}
