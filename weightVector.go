package ccsds

import "math"

func (ic *ImgCube) initializeDefaultWeights() []int {
	switch ic.p.PredictionMode {
	case FullPrediction:
		weights := make([]int, ic.p.NPreviousBands+3)
		if ic.p.NPreviousBands == 0 {
			return weights
		}
		weights[0] = 0
		weights[1] = 0
		weights[2] = 0
		weights[3] = 7 * (1<<ic.p.Omega - 3)
		for i := 4; i < ic.p.NPreviousBands; i++ {
			weights[i] = weights[i-1] >> 3
		}
		return weights
	case ReducedPrediction:
		weights := make([]int, ic.p.NPreviousBands)
		if ic.p.NPreviousBands == 0 {
			return weights
		}
		weights[0] = 7 * (1<<ic.p.Omega - 3)
		for i := 1; i < ic.p.NPreviousBands; i++ {
			weights[i] = weights[i-1] >> 3
		}
		return weights
	default:
		panic("Invalid prediction mode!")
	}
}

func (ic *ImgCube) initializeCustomWeights() [][]int {
	weights := make([][]int, ic.p.NZ*ic.p.Omega)
	for i := 0; i < ic.p.NZ-1; i++ {
		cZ := ic.p.Omega
		for j := 0; j < cZ; j++ {
			weights[i][j] = 2 << int(float64(ic.p.Omega+3-ic.p.Q)+math.Ceil(float64(2<<ic.p.Omega+2-ic.p.Q-1)))
		}
	}
	return weights
}

func (ic *ImgCube) InitializeWeights()
