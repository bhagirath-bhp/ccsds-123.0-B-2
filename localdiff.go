package ccsds_spike

import (
	"fmt"
)

func (ccsds *CCSDS) calcDiffVector(sum, z, y, x int) []int {
	var diffV []int
	pz := min(ccsds.config.predictionConfig.nPrevBands, z)
	switch ccsds.config.predictionConfig.pMode {
	case FullPredictionMode:
		diffV = make([]int, pz+3)
		dN, dW, dNW := ccsds.dirDiff(sum, z, y, x)
		diffV[0] = dN
		diffV[1] = dW
		diffV[2] = dNW
	case ReducedPredictionMode:
		diffV = make([]int, pz)
	default:
		panic(fmt.Errorf("Invalid prediction mode: %t", ccsds.config.predictionConfig.pMode))
	}
	for i := 0; i < pz; i++ {
		diffV = append(diffV, ccsds.centralDiff(z-i, y, x))
	}
	return diffV
}

func (ccsds *CCSDS) centralDiff(z, y, x int) int {
	return (4 * ccsds.srAt(z, y, x)) - ccsds.LocalSum(z, y, x)
}

func (ccsds *CCSDS) dirDiff(sum, z, y, x int) (int, int, int) {
	var dN, dW, dNW int
	dN = ccsds.dirDiffN(sum, z, y, x)
	dW = ccsds.dirDiffW(sum, z, y, x)
	dNW = ccsds.dirDiffNW(sum, z, y, x)
	return dN, dW, dNW
}

func (ccsds *CCSDS) dirDiffN(sum, z, y, x int) int {
	switch {
	case y > 0:
		return (4 * ccsds.srAt(z, y-1, x)) - sum
	case y == 0:
		return 0
	default:
		panic(fmt.Errorf("Invalid coordinates, unabled to calculate directionalLocaldifferenceN at x: %d, y: %d, z: %d!", x, y, z))
	}
}

func (ccsds *CCSDS) dirDiffW(sum, z, y, x int) int {
	switch {
	case x > 0 && y > 0:
		return (4 * ccsds.srAt(z, y, x-1)) - sum
	case x == 0 && y > 0:
		return (4 * ccsds.srAt(z, y-1, x)) - sum
	case y == 0:
		return 0
	default:
		panic(fmt.Errorf("Invalid coordinates, unabled to calculate directionalLocaldifferenceW at x: %d, y: %d, z: %d!", x, y, z))
	}
}

func (ccsds *CCSDS) dirDiffNW(sum, z, y, x int) int {
	switch {
	case x > 0 && y > 0:
		return (4 * ccsds.srAt(z, y-1, x-1)) - sum
	case x == 0 && y > 0:
		return (4 * ccsds.srAt(z, y-1, x)) - sum
	case y == 0:
		return 0
	default:
		panic(fmt.Errorf("Invalid coordinates, unabled to calculate directionalLocaldifferenceNW at x: %d, y: %d, z: %d!", x, y, z))
	}
}
