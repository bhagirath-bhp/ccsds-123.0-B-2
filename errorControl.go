package ccsds_spike

import "fmt"

func (ccsds *CCSDS) maxError(pred, z, y, x int) int {
	switch ccsds.config.genConfig.compressionMode {
	case Lossless:
		return 0
	case Lossy:
		aZ := ccsds.config.errorControl.aZ
		rZ := ccsds.config.errorControl.rZ
		switch {
		case aZ != 0 && rZ == 0:
			return ccsds.config.errorControl.aZ
		case rZ != 0 && aZ == 0:
			exp := 1 << ccsds.config.genConfig.dRange
			return ccsds.floorDiv(rZ*pred, exp)
		case rZ != 0 && aZ != 0:
			exp := 1 << ccsds.config.genConfig.dRange
			rErr := ccsds.floorDiv(rZ*pred, exp)
			return ccsds.min(aZ, rErr)
		default:
			panic(fmt.Errorf("Invalid error control inputs!"))
		}
	default:
		panic(fmt.Errorf("Invalid compression mode!"))
	}
}
