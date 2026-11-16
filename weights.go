package ccsds_spike

import "math"

func (ccsds *CCSDS) initializeWeights(z int) []int {
	var weights []int
	pz := min(ccsds.config.predictionConfig.nPrevBands, z)
	omega := ccsds.config.weightConfig.wCompRes
	switch ccsds.config.weightConfig.wInitMode {
	case DefaultWeightInit:
		idx := 0
		switch ccsds.config.predictionConfig.pMode {
		case FullPredictionMode:
			weights = make([]int, pz+3)
			weights[0] = 0
			weights[1] = 0
			weights[2] = 0
			idx = 3
		case ReducedPredictionMode:
			weights = make([]int, pz)
			idx = 0
		}
		wb0 := 7 * (1 << omega) / 8
		weights[idx] = wb0
		for i := 1; i < pz; i++ {
			wbi := weights[i+idx-1] / 8
			weights[i+idx] = wbi
		}
		return weights
	case CustomWeightInit:
		coff1 := 1 << (omega + 3 - ccsds.config.weightConfig.q)
		coffIdentity := (1 << (omega + 2 - ccsds.config.weightConfig.q)) - 1
		w_init_v_z := ccsds.config.weightConfig.wInitVec[z]
		weights := make([]int, len(w_init_v_z))
		for j := 0; j < len(w_init_v_z); j++ {
			weights[j] = (coff1 * w_init_v_z[j]) + coffIdentity
		}
		return weights

	default:
		panic("Invalid weight initialization mode!")
	}
}
func (ccsds *CCSDS) expOffsetIntraBAt(z int) int {
	if ccsds.config.weightConfig.expOffsetMode == false {
		return 0
	} else {
		return ccsds.config.weightConfig.expOffsetIntraBVec[z]
	}
}

func (ccsds *CCSDS) expOffsetInterBAt(z, i int) int {
	if ccsds.config.weightConfig.expOffsetMode == false {
		return 0
	} else {
		return ccsds.config.weightConfig.expOffsetInterBVec[z][i]
	}
}

func (ccsds *CCSDS) getUpdatedWT(z, y, x, qBinCenter, predDR, predError int, weightsVZ []int) []int {
	pZ := min(ccsds.config.predictionConfig.nPrevBands, z)
	predErrDR := 2*qBinCenter - predDR
	t := y*ccsds.image.nX + x
	omega := ccsds.config.weightConfig.wCompRes
	dRange := ccsds.config.genConfig.dRange
	tInc := ccsds.config.weightConfig.tInc
	nX := ccsds.image.nX
	var scalingExp int
	vMin := ccsds.config.weightConfig.vMin
	vMax := ccsds.config.weightConfig.vMax
	if t > 0 {
		scalingExp = ccsds.clip(vMin+ccsds.floorDiv(nX, tInc), vMin, vMax) + dRange - omega
	}
	e1 := weightsVZ[t]
	var resWeights []int
	idxOffset := 0
	resWeights = make([]int, pZ)
	if ccsds.config.predictionConfig.pMode == FullPredictionMode {
		resWeights = make([]int, pZ+3)
		e2n := math.Floor(1.0 / 2 * float64(ccsds.getSign_p(predErrDR)) * math.Exp2(float64(-1*(scalingExp+ccsds.expOffsetIntraBAt(z)))) * float64(ccsds.diffV[z][0]))
		e2w := math.Floor(1.0 / 2 * float64(ccsds.getSign_p(predErrDR)) * math.Exp2(float64(-1*(scalingExp+ccsds.expOffsetIntraBAt(z)))) * float64(ccsds.diffV[z][1]))
		e2nw := math.Floor(1.0 / 2 * float64(ccsds.getSign_p(predErrDR)) * math.Exp2(float64(-1*(scalingExp+ccsds.expOffsetIntraBAt(z)))) * float64(ccsds.diffV[z][2]))
		resWeights[0] = ccsds.clip(e1+int(e2n), ccsds.preCalculated.wMin, ccsds.preCalculated.wMax)
		resWeights[1] = ccsds.clip(e1+int(e2w), ccsds.preCalculated.wMin, ccsds.preCalculated.wMax)
		resWeights[2] = ccsds.clip(e1+int(e2nw), ccsds.preCalculated.wMin, ccsds.preCalculated.wMax)
		idxOffset = 3
	}
	for i := 0; i < pZ; i++ {
		e2 := math.Floor(1 / 2 * float64(ccsds.getSign_p(predErrDR)) * math.Exp2(float64(-1*(scalingExp+ccsds.expOffsetInterBAt(z, i)))) * float64(ccsds.diffV[z][i+idxOffset]))
		resWeights[i+idxOffset] = ccsds.clip(e1+int(e2), ccsds.preCalculated.wMin, ccsds.preCalculated.wMax)
	}
	return resWeights
}

/*
how should the weight embedding work in support both full and reduced prediction mode,
full prediction mode
[
1: N
2: W
3: NW
5. 1
6. 2
7. 3
8. 4
9. 5
10.6
11.7
12.8
13.9
14.10
]
*/

// and, when full prediction mode is used, for the directional components,
