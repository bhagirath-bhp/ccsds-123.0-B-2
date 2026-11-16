package ccsds_spike

import "fmt"

func (ccsds *CCSDS) predict(predDiff, sum, z, y, x int) (int, int, int) {
	omega := ccsds.config.weightConfig.wCompRes
	exp := 1 << omega
	v1 := ccsds.signedMod(predDiff+exp*(sum-4*ccsds.preCalculated.sMid), ccsds.config.genConfig.rSize)
	v2 := (exp * 4 * ccsds.preCalculated.sMid) + exp*2
	min := exp * 4 * ccsds.preCalculated.sMin
	max := (exp * 4 * ccsds.preCalculated.sMax) + exp*2
	predHR := ccsds.clip(v1+v2, min, max)
	pZ := ccsds.config.predictionConfig.nPrevBands
	var predDR, pred int
	switch {
	case x*y > 0:
		predDR = ccsds.floorDiv(predHR, exp*2)
	case (x*y == 0) && z > 0 && pZ > 0:
		predDR = ccsds.srAt(z-1, y, x)
	case (x*y == 0) && (pZ == 0 || z == 0):
		predDR = 2 * ccsds.preCalculated.sMid
	default:
		panic(fmt.Errorf("Bad coordinates found while calculating predDR!"))
	}
	pred = ccsds.floorDiv(predDR, 2)
	return pred, predDR, predHR
}
