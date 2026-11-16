package ccsds_spike

import "fmt"

func (ccsds *CCSDS) quantize(pred, sample, maxErr, z, y, x int) int {
	residual := sample - pred
	switch {
	case x*y == 0:
		return residual
	case x*y > 0:
		return ccsds.getSign(residual) * ccsds.floorDiv(residual+maxErr, 2*maxErr+1)
	default:
		panic(fmt.Errorf("Invalid coordinated encountered in quantization!: z,y,x: %d,%d,%d", z, y, x))
	}
}
