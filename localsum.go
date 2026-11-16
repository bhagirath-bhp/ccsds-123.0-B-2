package ccsds_spike

import "fmt"

func (ccsds *CCSDS) LocalSum(z, y, x int) int {
	// The value of σz,0,0 is not defined, as it is not needed.
	switch ccsds.config.predictionConfig.sumMode {
	case WideNeighborOriented:
		return ccsds.wideNeighborOrientedSum(z, y, x)
	case NarrowNeighborOriented:
		return ccsds.narrowNeighborOrientedSum(z, y, x)
	case WideColumnOriented:
		return ccsds.wideNeighborOrientedSum(z, y, x)
	case NarrowColumnOriented:
		return ccsds.narrowColumnOrientedSum(z, y, x)
	default:
		panic(fmt.Errorf("Invalid Operation: localSumMode, Unable to calculate localSum at x: %d, y: %d, z: %d!", x, y, z))
	}
}

func (ccsds *CCSDS) wideNeighborOrientedSum(z, y, x int) int {
	switch {
	case y > 0 && x > 0 && x < ccsds.image.nX-1:
		return (ccsds.srAt(z, y, x-1) +
			ccsds.srAt(z, y-1, x-1) +
			ccsds.srAt(z, y-1, x) +
			ccsds.srAt(z, y-1, x+1))
	case y == 0 && x > 0:
		return 4 * ccsds.srAt(z, y, x-1)
	case y > 0 && x == 0:
		return 2 * (ccsds.srAt(z, y-1, x) + ccsds.srAt(z, y-1, x+1))
	case y > 0 && x == ccsds.image.nX-1:
		return (ccsds.srAt(z, y, x-1) + ccsds.srAt(z, y-1, x-1) + ccsds.srAt(z, y-1, x))
	default:
		panic(fmt.Errorf("Invalid coordinates, unable to calculate WideNeighborOrientedSum at x: %d, y: %d, z: %d!", x, y, z))
	}
}

func (ccsds *CCSDS) narrowNeighborOrientedSum(z, y, x int) int {
	switch {
	case y > 0 && x > 0 && x < ccsds.image.nX-1:
		return (ccsds.srAt(z, y-1, x-1) + ccsds.srAt(z, y-1, x) + ccsds.srAt(z, y-1, x+1))
	case y == 0 && x > 0 && z > 0:
		return (4 * ccsds.srAt(z-1, y, x-1))
	case y > 0 && x == 0:
		return 2 * (ccsds.srAt(z, y-1, x) + ccsds.srAt(z, y-1, x+1))
	case y > 0 && x == ccsds.image.nX-1:
		return 2 * (ccsds.srAt(z, y-1, x-1) + ccsds.srAt(z, y-1, x))
	case y == 0 && x > 0 && z == 0:
		return (4 * ccsds.preCalculated.smid)
	default:
		panic(fmt.Errorf("Invalid coordinates, unable to calculate narrowNeighborOrientedSum at x: %d, y: %d, z: %d!", x, y, z))
	}
}

func (ccsds *CCSDS) wideColumnOrientedSum(z, y, x int) int {
	switch {
	case y > 0:
		return (4 * ccsds.srAt(z, y-1, x))
	case y == 0 && x > 0:
		return (4 * ccsds.srAt(z, y, x-1))
	default:
		panic(fmt.Errorf("Invalid coordinates, unable to calculate narrowColumnOrientedSum at x: %d, y: %d, z: %d!", x, y, z))
	}
}

func (ccsds *CCSDS) narrowColumnOrientedSum(z, y, x int) int {
	switch {
	case y > 0:
		return ccsds.srAt(z, y-1, x)
	case y == 0 && x > 0 && z > 0:
		return ccsds.srAt(z-1, y, x-1)
	case y == 0 && x > 0 && z == 0:
		return 4 * ccsds.preCalculated.smid
	default:
		panic(fmt.Errorf("Invalid coordinates, unable to calculate narrowColumnOrientedSum at x: %d, y: %d, z: %d!", x, y, z))
	}
}
