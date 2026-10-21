package ccsds

import "fmt"

func (ic *ImgCube) localDifferenceVector(z, y, x int) ([]int, error) {
	sumZYX := ic.localSum(z, y, x)
	switch ic.p.PredictionMode {
	case FullPrediction:
		diffV := make([]int, 3+ic.p.NPreviousBands)
		diffV[0] = ic.directionalLocalDifferenceN(z, y, x, sumZYX)
		diffV[1] = ic.directionalLocalDifferenceW(z, y, x, sumZYX)
		diffV[2] = ic.directionalLocalDifferenceNW(z, y, x, sumZYX)
		for i := 0; i < ic.p.NPreviousBands; i++ {
			diffV[i+2] = ic.centralLocalDifference(z, y, x, sumZYX)
		}
		return diffV, nil
	case ReducedPrediction:
		diffV := make([]int, ic.p.NPreviousBands)
		for i := 0; i < ic.p.NPreviousBands; i++ {
			diffV[i] = ic.centralLocalDifference(z, y, x, sumZYX)
		}
		return diffV, nil
	default:
		return nil, fmt.Errorf("Invalid prediction mode!")
	}
}

func (ic *ImgCube) centralLocalDifference(z, y, x, sumZYX int) int {
	sZYX := ic.sampleRepresentativeAt(z, y, x)
	return (4 * sZYX) - sumZYX
}

func (ic *ImgCube) directionalLocalDifferenceN(z, y, x, sumZYX int) int {
	switch {
	case y > 0:
		sZYm1X := ic.sampleRepresentativeAt(z, y-1, x)
		return (4 * sZYm1X) - sumZYX
	case y == 0:
		return 0
	default:
		panic("Invalid coords!")
	}
}

func (ic *ImgCube) directionalLocalDifferenceW(z, y, x, sumZYX int) int {
	switch {
	case x > 0 && y > 0:
		sZYXm1 := ic.sampleRepresentativeAt(z, y, x-1)
		return (4 * sZYXm1) - sumZYX
	case x == 0 && y > 0:
		sZYm1X := ic.sampleRepresentativeAt(z, y-1, x)
		return (4 * sZYm1X) - sumZYX
	case y == 0:
		return 0
	default:
		panic("Invalid coords!")
	}
}

func (ic *ImgCube) directionalLocalDifferenceNW(z, y, x, sumZYX int) int {
	switch {
	case x > 0 && y > 0:
		sZYm1Xm1 := ic.sampleRepresentativeAt(z, y-1, x-1)
		return (4 * sZYm1Xm1) - sumZYX
	case x == 0 && y > 0:
		sZYm1X := ic.sampleRepresentativeAt(z, y-1, x)
		return (4 * sZYm1X) - sumZYX
	case y == 0:
		return 0
	default:
		panic("Invalid coords!")
	}
}
