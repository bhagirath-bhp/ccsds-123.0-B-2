package ccsds_spike

import "fmt"

func (ccsds *CCSDS) srAt(z, y, x int) int {
	return ccsds.sr[(z*ccsds.image.nY*ccsds.image.nX)+
		(x*ccsds.image.nY)+
		x]
}

func (ccsds *CCSDS) getUpdatedSR(z, y, x, predHR, pred, quantized, maxError int) (int, int) {
	var srDR, sr int
	t := y*ccsds.image.nX + x
	res := ccsds.config.srConfig.srRes
	damping := ccsds.config.srConfig.srDamping
	offset := ccsds.config.srConfig.srOffset
	omega := ccsds.config.weightConfig.wCompRes
	exp := 1 << omega

	qBinCenter := ccsds.clip((pred+quantized)+(2*maxError+1), ccsds.preCalculated.sMin, ccsds.preCalculated.sMax)

	e1 := 4 * ((1 << res) - damping)
	e2 := qBinCenter*exp - ccsds.getSign(quantized)*maxError*offset*(1<<(omega-res))
	e3 := damping*predHR - damping*exp*2
	denominator := 1 << (omega*res + 1)

	srDR = ((e1 * e2) + e3) / denominator

	switch {
	case t == 0:
		sr = ccsds.srAt(0, 0, 0)
		return qBinCenter, sr
	case t > 0:
		sr = ccsds.floorDiv(srDR+1, 2)
		return qBinCenter, sr
	default:
		panic(fmt.Errorf("Invalid coordinated encountered in updateSR: z=%d, y=%d, x=%d", z, y, x))
	}
}
