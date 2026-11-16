package ccsds_spike

import (
	"fmt"
)

type CCSDS struct {
	config          Config
	image           Image
	preCalculated   preCalculated
	sr              []int
	mappedQuantized []int
	diffV           [][]int
}

type preCalculated struct {
	sMid      int
	sMin      int
	sMax      int
	imageSize int
	wMin      int
	wMax      int
}

func main() error {
	confPath := "userInput.toml"
	var ccsds CCSDS
	conf, err := ccsds.fetchConfig(confPath)
	if err != nil {
		return fmt.Errorf("Failed fetching config! %v", err)
	}
	ccsds.config = conf
	image, err := ccsds.fetchImage(ccsds.config.genConfig.imagePath)
	if err != nil {
		return fmt.Errorf("Failed fetching image! %v", err)
	}
	ccsds.image = image
	return nil
}

func (ccsds *CCSDS) preCalculate() {
	ccsds.preCalculated.imageSize = ccsds.image.nZ * ccsds.image.nY * ccsds.image.nX
	if ccsds.image.isSignedInt {
		exp := 1 << ccsds.config.genConfig.dRange
		ccsds.preCalculated.sMin = -1 * (exp / 2)
		ccsds.preCalculated.sMax = (exp / 2) - 1
		ccsds.preCalculated.sMid = 0
	} else {
		exp := 1 << ccsds.config.genConfig.dRange
		ccsds.preCalculated.sMin = 0
		ccsds.preCalculated.sMax = exp - 1
		ccsds.preCalculated.sMid = exp / 2
	}
	omega := ccsds.config.weightConfig.wCompRes
	ccsds.preCalculated.wMin = -1 * (1 << (omega + 2))
	ccsds.preCalculated.wMax = (1 << (omega + 2)) - 1
	ccsds.sr = make([]int, ccsds.preCalculated.imageSize)
}
func (ccsds *CCSDS) predError(pred, sample, z, y, x int) int {
	// TODO
	return 0
}
func (ccsds *CCSDS) mappedQuantize(quantized, pred, maxErr, z, y, x int) int {
	// TODO
	return 0
}
func (ccsds *CCSDS) ccsds1() {
	ccsds.validateInput(ccsds.config, ccsds.image)
	// initializeContext
	for z := 0; z < ccsds.image.nZ; z++ {
		weightsV := ccsds.initializeWeights(z)
		for y := 0; y < ccsds.image.nY; y++ {
			for x := 0; x < ccsds.image.nX; x++ {
				sample := ccsds.sampleAt(z, y, x)
				sum := ccsds.LocalSum(z, y, x)
				ccsds.diffV[z] = ccsds.calcDiffVector(sum, z, y, x)
				predDiff := ccsds.dotProduct(weightsV, ccsds.diffV[z])
				predHR, predDR, pred := ccsds.predict(predDiff, sum, z, y, x)
				maxErr := ccsds.maxError(pred, z, y, x)
				quantized := ccsds.quantize(pred, sample, maxErr, z, y, x)
				predErr := ccsds.predError(pred, sample, z, y, x)
				qBinCenter, updatedSR := ccsds.getUpdatedSR(z, y, x, predHR, pred, quantized, maxErr)
				ccsds.sr[(z*ccsds.image.nY*ccsds.image.nX)+
					(x*ccsds.image.nY)+
					x] = updatedSR
				ccsds.getUpdatedWT(z, y, x, qBinCenter, predDR, predErr, weightsV)
				mappedQuantized := ccsds.mappedQuantize(quantized, pred, maxErr, z, y, x)
				idx := z * y * x // temp index, TODO: update with real idx calculation
				ccsds.mappedQuantized[idx] = mappedQuantized
			}
		}
	}

}

// precompute list
/*
	smin
	smax
	weight update scaling exponent p(t)

*/
